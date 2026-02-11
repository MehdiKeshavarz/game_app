package accesscontrol

import (
	"fmt"
	"game_app/entity"
	"game_app/pkg/errmsg"
	"game_app/pkg/richerror"
	"game_app/pkg/slice"
	"game_app/repository/mysql"
	"strings"
)

func (d *DB) GetUserPermissionsTitles(userID uint, role entity.Role) ([]entity.PermissionTitle, error) {
	const op = "mysql.GetUserPermissionsTitles"

	roleACL := make([]entity.AccessControl, 0)

	rows, err := d.conn.Conn().Query(`select * from access_controls where actor_type = ?  and actor_id = ? `, entity.RoleActorType, role)
	if err != nil {

		return nil, richerror.New(op).
			SetWrappedError(err).
			SetMessage(errmsg.ErrorSomethingWentWrong).
			SetKind(richerror.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		acl, err := scanAccessControl(rows)
		if err != nil {

			return nil, richerror.New(op).
				SetWrappedError(err).
				SetMessage(errmsg.ErrorSomethingWentWrong).
				SetKind(richerror.KindUnexpected)
		}

		roleACL = append(roleACL, acl)

	}

	fmt.Println("roleAcllllll:", roleACL)

	if err := rows.Err(); err != nil {
		return nil, richerror.New(op).
			SetWrappedError(err).
			SetMessage(errmsg.ErrorSomethingWentWrong).
			SetKind(richerror.KindUnexpected)
	}

	userACL := make([]entity.AccessControl, 0)

	fmt.Println("beforuserAcl: ", userACL)
	userRows, err := d.conn.Conn().Query(`select * from access_controls where actor_type = ?  and actor_id = ? `, entity.UserActorType, userID)
	if err != nil {
		return nil, richerror.New(op).
			SetWrappedError(err).
			SetMessage(errmsg.ErrorSomethingWentWrong).
			SetKind(richerror.KindUnexpected)
	}

	defer userRows.Close()

	for userRows.Next() {
		acl, err := scanAccessControl(userRows)
		if err != nil {
			fmt.Println("eerrr ", err)
			return nil, richerror.New(op).
				SetWrappedError(err).
				SetMessage(errmsg.ErrorSomethingWentWrong).
				SetKind(richerror.KindUnexpected)
		}

		userACL = append(userACL, acl)

	}
	fmt.Println("afteruserAcl: ", userACL)

	if err := userRows.Err(); err != nil {
		fmt.Println("eerrr ", err)
		return nil, richerror.New(op).
			SetWrappedError(err).
			SetMessage(errmsg.ErrorSomethingWentWrong).
			SetKind(richerror.KindUnexpected)
	}

	// merge acls by permission id

	PermissionIDs := make([]uint, 0)
	fmt.Println("beforPermissionID: ", PermissionIDs)
	for _, r := range roleACL {
		if !slice.DoesExist(PermissionIDs, r.PermissionID) {
			PermissionIDs = append(PermissionIDs, r.PermissionID)
		}

	}
	fmt.Println("afterPermissionID: ", PermissionIDs)

	if len(PermissionIDs) == 0 {
		return nil, nil
	}

	args := make([]any, len(PermissionIDs))
	fmt.Println("beforArgs: ", args)

	for i, id := range PermissionIDs {
		args[i] = id
	}
	fmt.Println("AfterArgs: ", args)
	query := "select * from permissions where id in (?" +
		strings.Repeat(",?", len(PermissionIDs)-1) + ")"

	pRows, err := d.conn.Conn().Query(query, args...)

	fmt.Println("pRows:", pRows)
	if err != nil {
		fmt.Println("eerrr ", err)
		return nil, richerror.New(op).
			SetWrappedError(err).
			SetMessage(errmsg.ErrorSomethingWentWrong).
			SetKind(richerror.KindUnexpected)
	}

	defer pRows.Close()

	permissionsTitle := make([]entity.PermissionTitle, 0)
	fmt.Println("beforPermissionsTitle: ", permissionsTitle)
	for pRows.Next() {
		fmt.Println("rows.next")
		per, err := scanPermission(pRows)
		if err != nil {
			fmt.Println("eerrr ", err)
			return nil, richerror.New(op).
				SetWrappedError(err).
				SetMessage(errmsg.ErrorSomethingWentWrong).
				SetKind(richerror.KindUnexpected)
		}

		permissionsTitle = append(permissionsTitle, per.Title)
	}

	fmt.Println("AfterPermissionsTitle: ", permissionsTitle)
	if err := userRows.Err(); err != nil {
		fmt.Println("eerrr ", err)
		return nil, richerror.New(op).
			SetWrappedError(err).
			SetMessage(errmsg.ErrorSomethingWentWrong).
			SetKind(richerror.KindUnexpected)
	}

	fmt.Println("permissionsTitle", permissionsTitle)
	return permissionsTitle, nil

}

func scanAccessControl(scanner mysql.Scanner) (entity.AccessControl, error) {
	var acl entity.AccessControl
	var createdAt []uint8
	err := scanner.Scan(&acl.ID, &acl.ActorID, &acl.ActorType, &acl.PermissionID, &createdAt)

	return acl, err

}

func scanPermission(scanner mysql.Scanner) (entity.Permission, error) {
	var per entity.Permission
	var createdAt []uint8
	err := scanner.Scan(&per.ID, &per.Title, &createdAt)

	return per, err
}

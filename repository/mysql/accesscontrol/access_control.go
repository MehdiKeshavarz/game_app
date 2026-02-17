package accesscontrol

import (
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

	if err := rows.Err(); err != nil {
		return nil, richerror.New(op).
			SetWrappedError(err).
			SetMessage(errmsg.ErrorSomethingWentWrong).
			SetKind(richerror.KindUnexpected)
	}

	userACL := make([]entity.AccessControl, 0)

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
			return nil, richerror.New(op).
				SetWrappedError(err).
				SetMessage(errmsg.ErrorSomethingWentWrong).
				SetKind(richerror.KindUnexpected)
		}

		userACL = append(userACL, acl)

	}
	if err := userRows.Err(); err != nil {
		return nil, richerror.New(op).
			SetWrappedError(err).
			SetMessage(errmsg.ErrorSomethingWentWrong).
			SetKind(richerror.KindUnexpected)
	}

	PermissionIDs := make([]uint, 0)
	for _, r := range roleACL {
		if !slice.DoesExist(PermissionIDs, r.PermissionID) {
			PermissionIDs = append(PermissionIDs, r.PermissionID)
		}

	}

	if len(PermissionIDs) == 0 {
		return nil, nil
	}

	args := make([]any, len(PermissionIDs))

	for i, id := range PermissionIDs {
		args[i] = id
	}
	query := "select * from permissions where id in (?" +
		strings.Repeat(",?", len(PermissionIDs)-1) + ")"

	pRows, err := d.conn.Conn().Query(query, args...)

	if err != nil {
		return nil, richerror.New(op).
			SetWrappedError(err).
			SetMessage(errmsg.ErrorSomethingWentWrong).
			SetKind(richerror.KindUnexpected)
	}

	defer pRows.Close()

	permissionsTitle := make([]entity.PermissionTitle, 0)
	for pRows.Next() {
		per, err := scanPermission(pRows)
		if err != nil {

			return nil, richerror.New(op).
				SetWrappedError(err).
				SetMessage(errmsg.ErrorSomethingWentWrong).
				SetKind(richerror.KindUnexpected)
		}

		permissionsTitle = append(permissionsTitle, per.Title)
	}

	if err := userRows.Err(); err != nil {
		return nil, richerror.New(op).
			SetWrappedError(err).
			SetMessage(errmsg.ErrorSomethingWentWrong).
			SetKind(richerror.KindUnexpected)
	}

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

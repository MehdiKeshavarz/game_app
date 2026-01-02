package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"game_app/entity"
	"game_app/pkg/errmsg"
	"game_app/pkg/richerror"
)

func (d Mysqldb) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"
	row := d.db.QueryRow(`select * from users where phone_number = ? `, phoneNumber)

	_, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}

		return false, richerror.New(op).
			SetWrappedError(err).
			SetMessage(errmsg.ErrorMsgCantScanQueryResult).
			SetKind(richerror.KindUnexpected)
	}

	return false, nil

}

func (d Mysqldb) Register(user entity.User) (entity.User, error) {
	res, err := d.db.Exec(`insert into users(name, phone_number,password) values(?, ?, ?)`, user.Name, user.PhoneNumber, user.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("can't execute command: %w", err)
	}

	// error is always nil
	id, _ := res.LastInsertId()
	user.ID = uint(id)

	return user, nil
}

func (d Mysqldb) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	const op = "mysql.GetUserByPhoneNumber"
	row := d.db.QueryRow(`select * from users where phone_number = ? `, phoneNumber)

	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, richerror.New(op).
				SetWrappedError(err).
				SetMessage(errmsg.ErrorMsgNotFound).
				SetKind(richerror.KindNotFound)
		}

		return entity.User{}, richerror.New(op).
			SetWrappedError(err).
			SetMessage(errmsg.ErrorMsgCantScanQueryResult).
			SetKind(richerror.KindUnexpected)
	}

	return user, nil

}

func (d Mysqldb) GetUserByID(userID uint) (entity.User, error) {
	const op = "mysql.GetUserByID"
	row := d.db.QueryRow(`select * from users where id = ? `, userID)

	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, richerror.New(op).
				SetWrappedError(err).
				SetMessage(errmsg.ErrorMsgNotFound).
				SetKind(richerror.KindNotFound)
		}
		return entity.User{}, richerror.New(op).
			SetWrappedError(err).
			SetMessage(errmsg.ErrorMsgCantScanQueryResult).
			SetKind(richerror.KindUnexpected)
	}

	return user, nil
}

func scanUser(row *sql.Row) (entity.User, error) {
	var user entity.User
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password)

	return user, err
}

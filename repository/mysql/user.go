package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"game_app/entity"
)

func (d Mysqldb) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	row := d.db.QueryRow(`select * from users where phone_number = ? `, phoneNumber)
	user := entity.User{}

	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil
		}

		return false, fmt.Errorf("can't scan query result: %w", err)
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

func (d Mysqldb) GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error) {

	row := d.db.QueryRow(`select * from users where phone_number = ? `, phoneNumber)
	user := entity.User{}

	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, false, nil
		}

		return entity.User{}, false, fmt.Errorf("can't scan query result: %w", err)
	}

	return user, true, nil

}

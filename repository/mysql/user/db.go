package user

import "game_app/repository/mysql"

type DB struct {
	conn *mysql.Mysqldb
}

func New(conn *mysql.Mysqldb) *DB {
	return &DB{
		conn: conn,
	}
}

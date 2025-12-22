package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     int
	DBName   string
}

type Mysqldb struct {
	config Config
	db     *sql.DB
}

func New(config Config) *Mysqldb {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s", config.Username,
		config.Password, config.Host, config.Port, config.DBName))
	if err != nil {
		panic(fmt.Errorf("can't open db connection: %v", err))
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &Mysqldb{db: db, config: config}
}

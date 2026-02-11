package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	DBName   string `koanf:"db_name"`
}

type Mysqldb struct {
	config Config
	db     *sql.DB
}

func (m *Mysqldb) Conn() *sql.DB {
	return m.db
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

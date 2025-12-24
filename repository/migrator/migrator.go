package migrator

import (
	"database/sql"
	"fmt"
	"game_app/repository/mysql"

	"github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect    string
	dbConfig   mysql.Config
	migrations *migrate.FileMigrationSource
}

func New(dbConfig mysql.Config, dialect string) Migrator {
	migrations := &migrate.FileMigrationSource{
		Dir: "./repository/mysql/migrations",
	}

	return Migrator{migrations: migrations, dbConfig: dbConfig, dialect: dialect}
}

func (m Migrator) Up() {
	db, err := sql.Open(m.dialect, fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true ",
		m.dbConfig.Username,
		m.dbConfig.Password,
		m.dbConfig.Host,
		m.dbConfig.Port,
		m.dbConfig.DBName))
	if err != nil {
		panic(fmt.Errorf("can't open db connection: %v", err))
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		panic(fmt.Errorf("can't apply migrations: %v", err))
	}
	fmt.Printf("Applied %d migrations!\n", n)
}

func (m Migrator) Down() {
	db, err := sql.Open(m.dialect, fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		m.dbConfig.Username,
		m.dbConfig.Password,
		m.dbConfig.Host,
		m.dbConfig.Port,
		m.dbConfig.DBName))
	if err != nil {
		panic(fmt.Errorf("can't open db connection: %v", err))
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		panic(fmt.Errorf("can't rollback migrations: %v", err))
	}
	fmt.Printf("Rollback %d migrations!\n", n)

}

func (m Migrator) Status() {
	// TODO - add status
}

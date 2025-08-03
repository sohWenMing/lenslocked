package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type DBConnections struct {
	UserService *UserService
	DB          *sql.DB
}

type pgConfig struct {
	host, port, user, password, dbname, sslmode string
}

func (p pgConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.host, p.port, p.user, p.password, p.dbname, p.sslmode,
	)
}

var config = pgConfig{
	"localhost",
	"5432",
	"baloo",
	"junglebook",
	"lenslocked",
	"disable",
}

func InitDBConnections() *DBConnections {
	db, err := sql.Open(
		"pgx",
		config.String(),
	)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	userServicePtr := &UserService{
		db,
	}
	fmt.Println("DB Connection has been initialised")
	dbc := &DBConnections{
		userServicePtr,
		db,
	}
	dbc.InitCreatedTablesIfNotExist()
	return dbc
}

func (dbc *DBConnections) InitCreatedTablesIfNotExist() {
	dbc.UserService.CreateUserTableIfNotExist()
}

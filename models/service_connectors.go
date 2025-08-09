package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type DBConnections struct {
	UserService    *UserService
	SessionService *SessionService
	DB             *sql.DB
}

type pgConfig struct {
	host, port, user, password, dbname, sslmode string
}

func (p pgConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.host, p.port, p.user, p.password, p.dbname, p.sslmode,
	)
}

func defaultConfig() pgConfig {
	return pgConfig{
		"localhost",
		"5432",
		"baloo",
		"junglebook",
		"lenslocked",
		"disable",
	}
}

func InitDBConnections() (dbc *DBConnections, err error) {
	db, err := sql.Open(
		"pgx",
		defaultConfig().String(),
	)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	sessionServicePtr := &SessionService{
		db,
	}
	userServicePtr := &UserService{
		db,
		sessionServicePtr,
	}
	fmt.Println("DB Connection has been initialised")
	dbc = &DBConnections{
		userServicePtr,
		sessionServicePtr,
		db,
	}
	return dbc, nil
}

package models

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

type DBConnections struct {
	UserService     *UserService
	SessionService  *SessionService
	ForgotPWService *ForgotPWService
	GalleryService  *GalleryService
	DB              *sql.DB
}

type PgConfig struct {
	host, port, user, password, dbname, sslmode string
}

func (p PgConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.host, p.port, p.user, p.password, p.dbname, p.sslmode,
	)
}
func (p PgConfig) DBInterface() {
}

func DefaultConfig() PgConfig {
	return PgConfig{
		"localhost",
		"5432",
		"baloo",
		"junglebook",
		"lenslocked",
		"disable",
	}
}

type DBConfig interface {
	String() string
	DBInterface()
}

func InitDBConnections(config DBConfig) (dbc *DBConnections, err error) {
	fmt.Println("eval connection string: ", config.String())
	db, err := sql.Open(
		"pgx",
		config.String(),
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
	forgotEmailServicePtr := &ForgotPWService{
		db,
	}
	galleryServicePtr := &GalleryService{
		db, "",
	}
	fmt.Println("DB Connection has been initialised")
	dbc = &DBConnections{
		userServicePtr,
		sessionServicePtr,
		forgotEmailServicePtr,
		galleryServicePtr,
		db,
	}
	return dbc, nil
}

func Migrate(db *sql.DB, dir string, embedMigrations embed.FS) error {
	goose.SetBaseFS(embedMigrations)
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

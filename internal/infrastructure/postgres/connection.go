package postgres

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ConnectionConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func (c ConnectionConfig) dsn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Database,
	)
}

func (c ConnectionConfig) url() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.Database,
	)
}

func NewConnection(config ConnectionConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.dsn()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	return db, nil
}

func RunMigrations(config ConnectionConfig) error {
	m, err := migrate.New("file://migrations", config.url())
	if err != nil {
		return fmt.Errorf("initializing migrations: %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("running migrations: %w", err)
	}
	return nil
}

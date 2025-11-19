package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	// needed for file source driver in migrate.
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

type DBDeps struct {
	DSN            string
	MigrationsPath string
}

func New(deps DBDeps) (*Storage, error) {
	storage := &Storage{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := storage.Connect(ctx, deps.DSN, deps.MigrationsPath)
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func (s *Storage) Connect(ctx context.Context, dsn string, migrationsPath string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return errors.Join(errors.New("failed open db"), err)
	}
	err = runMigrations(db, migrationsPath)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)
	if db == nil {
		return errors.New("db is not init")
	}
	err = db.PingContext(ctx)
	if err != nil {
		return errors.Join(errors.New("failed ping context db"), err)
	}
	s.db = db
	return nil
}

func runMigrations(db *sql.DB, migrationsPath string) error {
	if migrationsPath == "" {
		return nil
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Join(errors.New("failed create driver for migrations"), err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///"+migrationsPath,
		"postgres", driver)
	if err != nil {
		return errors.Join(errors.New("failed create migrate"), err)
	}
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return errors.Join(errors.New("failed up migrate"), err)
	}
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

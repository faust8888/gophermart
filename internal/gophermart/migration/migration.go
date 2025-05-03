package migration

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/faust8888/gophermart/internal/gophermart/repository/postgres"
	"github.com/golang-migrate/migrate/v4"
	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func Run(dataSourceName string) error {
	db, err := sql.Open(postgres.DatabaseDriverName, dataSourceName)
	if err != nil {
		return fmt.Errorf("migration.run: oppening conection - %w", err)
	}
	defer db.Close()

	driver, err := postgresMigrate.WithInstance(db, &postgresMigrate.Config{})
	if err != nil {
		return fmt.Errorf("migration.run: creating migration driver - %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/gophermart/migration/sql",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("migration.run: creating migration instance - %w", err)
	}

	if err = m.Up(); !errors.Is(err, migrate.ErrNoChange) && err != nil {
		return fmt.Errorf("migration.run: applying migrations - %w", err)
	}
	log.Printf("migration.run: migration applied successfully")
	return nil
}

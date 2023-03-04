package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type dataDB struct {
	databaseDNS string
}

func NewDB(databaseDNS string) dataDB {
	db, err := sql.Open("postgres", databaseDNS)

	if err != nil {
		log.Fatal("unable to open db")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		log.Fatalf("unable to init db driver %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://schema",
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatalf("unable to init db migrator %v", err)
	}

	m.Up()

	return dataDB{databaseDNS: databaseDNS}
}

func (db dataDB) Get(shortenURL string) (string, error) {
	return "", fmt.Errorf("todo: доделать")
}

func (db dataDB) GetItemsOfUser(userID string) ([]item, error) {
	return nil, fmt.Errorf("todo: доделать")
}

func (db dataDB) Add(originalURL, shortenURL, userID string) bool {
	return false
}

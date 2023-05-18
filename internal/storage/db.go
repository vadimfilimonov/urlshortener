package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/VadimFilimonov/urlshortener/internal/constants"
	utils "github.com/VadimFilimonov/urlshortener/internal/utils/generateid"
)

type dataDB struct {
	databaseDNS string
}

func RunMigrations(databaseDNS string) error {
	db, err := sql.Open("postgres", databaseDNS)

	if err != nil {
		db.Close()
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		db.Close()
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://schema",
		"postgres",
		driver,
	)
	if err != nil {
		db.Close()
		return err
	}

	m.Up()
	fmt.Println("Migrations completed!")
	return db.Close()
}

func NewDB(databaseDNS string) dataDB {
	return dataDB{databaseDNS: databaseDNS}
}

func (data dataDB) Get(shortenURL string) (string, error) {
	db, err := sql.Open("postgres", data.databaseDNS)

	if err != nil {
		db.Close()
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var originalURL string
	err = db.QueryRowContext(ctx, "SELECT original_url FROM urls WHERE shorten_url = $1 LIMIT 1", shortenURL).Scan(&originalURL)

	if err != nil {
		db.Close()
		return "", err
	}

	db.Close()
	return originalURL, nil
}

func (data dataDB) GetItemsOfUser(userID string) ([]item, error) {
	items := make([]item, 0)
	db, err := sql.Open("postgres", data.databaseDNS)

	if err != nil {
		db.Close()
		log.Println(err)
		return items, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, "SELECT * FROM urls WHERE user_id = $1", userID)

	if err != nil {
		db.Close()
		log.Println(err)
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var item item
		err = rows.Scan(&id, &item.userID, &item.ShortenURL, &item.OriginalURL)

		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	db.Close()
	return items, nil
}

func (data dataDB) Add(originalURL, userID string) (string, error) {
	fmt.Println("Start of Add")
	db, err := sql.Open("postgres", data.databaseDNS)

	if err != nil {
		db.Close()
		return "", err
	}

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO urls(user_id, shorten_url, original_url, status) VALUES($1,$2,$3,$4) ON CONFLICT (original_url) DO NOTHING")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	shortenURLPath := utils.GenerateID()

	sqlResult, err := stmt.ExecContext(ctx, userID, shortenURLPath, originalURL, itemStatusCreated)
	if err != nil {
		db.Close()
		return "", err
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		db.Close()
		return "", err
	}

	err = tx.Commit()

	if err != nil {
		return "", err
	}

	hasURLBeenAdded := rowsAffected != 0
	if !hasURLBeenAdded {
		err = db.QueryRowContext(ctx, "SELECT shorten_url FROM urls WHERE original_url = $1 LIMIT 1", originalURL).Scan(&shortenURLPath)

		if err != nil {
			db.Close()
			return "", err
		}

		return shortenURLPath, constants.ErrURLAlreadyExists
	}

	db.Close()
	return shortenURLPath, nil
}

func (data dataDB) Delete(ids []string, userID string) error {
	db, err := sql.Open("postgres", data.databaseDNS)

	if err != nil {
		db.Close()
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, "UPDATE urls SET status = $1 WHERE user_id = $2 and shorten_url IN ($3)", itemStatusDeleted, userID, strings.Join(ids, ","))
	db.Close()

	return err
}

package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"

	"github.com/VadimFilimonov/urlshortener/internal/constants"
	utils "github.com/VadimFilimonov/urlshortener/internal/utils/generateid"
)

type dataDB struct {
	db *sql.DB
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://schema",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	m.Up()
	return nil
}

func InitDB(databaseDNS string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseDNS)

	if err != nil {
		db.Close()
		return nil, err
	}

	err = runMigrations(db)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func NewDB(db *sql.DB) dataDB {
	return dataDB{db: db}
}

func (data dataDB) Get(shortenURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var originalURL string
	var status string
	err := data.db.QueryRowContext(ctx, "SELECT original_url, status FROM urls WHERE shorten_url = $1 LIMIT 1", shortenURL).Scan(&originalURL, &status)

	if err != nil {
		return "", err
	}

	if status == itemStatusDeleted {
		return "", ErrURLHasBeenDeleted
	}

	return originalURL, nil
}

func (data dataDB) GetItemsOfUser(userID string) ([]item, error) {
	items := make([]item, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := data.db.QueryContext(ctx, "SELECT * FROM urls WHERE user_id = $1", userID)

	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var item item
		err = rows.Scan(&id, &item.userID, &item.ShortenURL, &item.OriginalURL, &item.status)

		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (data dataDB) Add(originalURL, userID string) (string, error) {
	tx, err := data.db.Begin()
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
		return "", err
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return "", err
	}

	err = tx.Commit()

	if err != nil {
		return "", err
	}

	hasURLBeenAdded := rowsAffected != 0
	if !hasURLBeenAdded {
		err = data.db.QueryRowContext(ctx, "SELECT shorten_url FROM urls WHERE original_url = $1 LIMIT 1", originalURL).Scan(&shortenURLPath)

		if err != nil {
			return "", err
		}

		return shortenURLPath, constants.ErrURLAlreadyExists
	}

	return shortenURLPath, nil
}

func (data dataDB) Delete(ids []string, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE urls SET status = $1 WHERE user_id = $2 and shorten_url = ANY($3)"
	_, err := data.db.ExecContext(ctx, query, itemStatusDeleted, userID, pq.Array(ids))
	return err
}

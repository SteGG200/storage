package db

import (
	"database/sql"
	"errors"

	"github.com/SteGG200/storage/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sqlx.DB
}

func New(dbFilePath string) (db *DB, err error) {
	connection, err := sqlx.Connect("sqlite3", dbFilePath)

	if err != nil {
		logger.ErrorLogger.Fatal(err)
		return nil, err
	}

	db = &DB{
		connection,
	}

	logger.InfoLogger.Print("Connected to database successfully")

	return
}

func (db *DB) GetVersion() (version string, err error) {
	query := "SELECT version FROM version"
	row := db.QueryRow(query)
	if err := row.Scan(&version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		logger.ErrorLogger.Print(err)
		return "", err
	}

	return version, nil
}

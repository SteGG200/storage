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

/*
New creates a new database connection.

Params:

	dbFilePath string // The path to the SQLite database file.

Returns:

	*DB // The new database connection.
	error // An error if any occurred during the database connection process.
*/
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

/*
GetVersion retrieves the current version of the database.

Returns:

	string // The current version of the database.
	error // An error if any occurred during the version retrieval process.
*/
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

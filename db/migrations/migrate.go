package main

import (
	"embed"
	"errors"
	"flag"
	"io/fs"
	"sort"
	"strings"

	"github.com/SteGG200/storage/db"
	"github.com/SteGG200/storage/logger"
)

var dbFilePath = flag.String("database", "./data.db", "Path to your database file.")

//go:embed *.sql
var sqlFiles embed.FS

func main() {
	flag.Parse()

	db, err := db.New(*dbFilePath)

	if err != nil {
		logger.ErrorLogger.Fatal(err)
		return
	}

	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS version (version VARCHAR NOT NULL)"); err != nil {
		logger.ErrorLogger.Fatal(err)
		return
	}

	currentVersion, err := db.GetVersion()

	if err != nil {
		logger.ErrorLogger.Fatal(err)
		return
	}

	versions, err := getSchemaFiles()

	if err != nil {
		logger.ErrorLogger.Fatal(err)
		return
	}

	if len(versions) == 0 {
		logger.WarningLogger.Println("No new schema files found.")
		return
	}

	var newVersion string

	for _, version := range versions {
		if currentVersion == "" || version > currentVersion {
			sqlQuery, err := fs.ReadFile(sqlFiles, version+".sql")

			if err != nil {
				logger.ErrorLogger.Print("Faile to apply schema %s.sql", version)
				logger.ErrorLogger.Fatal(err)
				return
			}

			tx, err := db.Begin()

			if _, err = tx.Exec(string(sqlQuery)); err != nil {
				tx.Rollback()
				logger.ErrorLogger.Print("Failed to apply schema %s.sql", version)
				logger.ErrorLogger.Fatal(err)
				return
			}

			tx.Commit()

			logger.InfoLogger.Printf("Applied schema %s.sql to database successfully", version)
			newVersion = version
		}
	}

	if newVersion != currentVersion {
		if currentVersion == "" {
			if _, err := db.Exec("INSERT INTO version(version) VALUES (?)", newVersion); err != nil {
				logger.ErrorLogger.Fatal(err)
			}
		} else {
			if _, err := db.Exec("UPDATE version SET version = ?", newVersion); err != nil {
				logger.ErrorLogger.Fatal(err)
			}
		}
	} else {
		logger.WarningLogger.Print("Database is already updated to latest version")
	}
}

func getSchemaFiles() (versions []string, err error) {
	files, err := fs.ReadDir(sqlFiles, ".")

	if err != nil {
		return nil, err
	}

	versions = make([]string, 0)

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			return nil, errors.New("Invalid schema file: " + file.Name())
		}
		version := strings.TrimSuffix(file.Name(), ".sql")
		versions = append(versions, version)
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i] < versions[j]
	})

	return
}

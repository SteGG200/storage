package db

import (
	"database/sql"
	"errors"

	"github.com/SteGG200/storage/logger"
)

func SaveUploadSession(db *DB, session map[string]any) error {
	query := `
	INSERT INTO upload_session 
	VALUES(?, ?, ?, ?)
	`
	_, err := db.Exec(query, session["token"], session["path"], session["filename"], session["createdAt"])

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	return nil
}

func GetUploadSession(db *DB, tokenString string) (result *UploadSession, err error) {
	query := `
	SELECT path, filename, created_at 
	FROM upload_session 
	WHERE token = ?
	`
	result = &UploadSession{}
	err = db.Get(result, query, tokenString)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		logger.ErrorLogger.Print(err)
		return nil, err
	}

	return
}

func RemoveUploadSessionByToken(db *DB, tokenString string) error {
	query := `
	DELETE FROM upload_session 
	WHERE token =?
	`
	_, err := db.Exec(query, tokenString)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	return nil
}

func RemoveUploadSessionByPath(db *DB, path string, filename string) error {
	query := `
	DELETE FROM upload_session 
	WHERE path = ? 
	AND filename = ?
	`
	_, err := db.Exec(query, path, filename)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	return nil
}

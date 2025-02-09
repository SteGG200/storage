package db

import (
	"database/sql"
	"errors"

	"github.com/SteGG200/storage/logger"
)

/*
SaveUploadSession saves an upload session into the database.
*/
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

/*
GetUploadSession retrieves an upload session by its token from the database.
If no session is found, it returns nil and nil.
*/
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

/*
RemoveUploadSessionByToken removes an upload session by its token from the database.
*/
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

/*
RemoveUploadSessionByPath removes all upload sessions by their path and filename from the database.
*/
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

package db

import "github.com/SteGG200/storage/logger"

func CheckIfPathHasAuth(db *DB, path string) (bool, error) {
	query := `
  SELECT id 
  FROM item_password 
  WHERE path = ?
  `

	rows, err := db.Query(query, path)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return false, err
	}

	defer rows.Close()

	return rows.Next(), nil
}

func CreateAuthForPath(db *DB, path string, password string) error {
	query := `
	INSERT INTO item_password (path, password)
	VALUES (?,?)
	`

	_, err := db.Exec(query, path, password)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	return nil
}

func GetPasswordOfPath(db *DB, path string) (string, error) {
	query := `
	SELECT password
	FROM item_password
	WHERE path =?
	`

	var row ItemPassword
	err := db.Get(&row, query, path)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return "", err
	}

	return row.Password, nil
}

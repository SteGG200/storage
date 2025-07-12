package db

import "github.com/SteGG200/storage/logger"

/*
CheckIfPathHasAuth checks if a given path requires authentication to access	.
*/
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

/*
CreateAuthForPath inserts a new path with the given password into the database
*/
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

func ChangePasswordOfPath(db *DB, path string, password string) error {
	query := `
	UPDATE item_password
	SET password = ?
	WHERE path = ?
	`

	_, err := db.Exec(query, password, path)

	if err != nil {
		logger.ErrorLogger.Print(err)
		return err
	}

	return nil
}

/*
GetPasswordOfPath retrieves the password for a given path from the database.
*/
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

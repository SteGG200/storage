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

package utils

import (
	"github.com/SteGG200/storage/db"
	"github.com/SteGG200/storage/logger"
	"github.com/SteGG200/storage/server/exception"
)

/*
VerifyUploadSession verifies that the provided upload session token is valid.

Params:

	database *db.DB // The database connection.
	tokenString string // The upload session token to verify.

Returns:

	*db.UploadSession // The validated upload session, or nil if the token is invalid.
	error // An error if any occurred during the session verification process.
*/
func VerifyUploadSession(database *db.DB, tokenString string) (session *db.UploadSession, err error) {
	data, err := ParseToken(tokenString)

	if err != nil {
		return nil, err
	}

	session, err = db.GetUploadSession(database, tokenString)

	if err != nil {
		return nil, err
	}

	if session == nil || data["filename"] != session.Filename || data["path"] != session.Path || int64(data["createdAt"].(float64)) != session.CreatedAt {
		logger.ErrorLogger.Print(exception.ErrInvalidToken)
		return nil, exception.ErrInvalidToken
	}

	return
}

package utils

import (
	"github.com/SteGG200/storage/db"
	"github.com/SteGG200/storage/logger"
	"github.com/SteGG200/storage/server/exception"
)

func Verify(database *db.DB, tokenString string) (session *db.UploadSession, err error) {
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

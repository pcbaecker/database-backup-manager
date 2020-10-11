package readdb

import (
	"os"
	"errors"
	"database-backup-manager/constants"
)

type ReaddbService interface {
	ReadDatabaseNames() ([]string,error)
}

func NewReaddbService() (ReaddbService,error) {
	db_type := os.Getenv(constants.DATABASE_TYPE)
	if (db_type == "mysql") {
		return NewMysqlReaddbService()
	}

	return nil, errors.New("Unknown database type = " + db_type)
}
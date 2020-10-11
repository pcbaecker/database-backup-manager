package backup

import (
	"os"
	"errors"
	"database-backup-manager/constants"
)

type DatabaseBackupService interface {
	BackupDatabases(databaseNames []string) (string,bool,error)
}

func NewDatabaseBackupService() (DatabaseBackupService,error) {
	database_type := os.Getenv(constants.DATABASE_TYPE)
	if (database_type == "mysql") {
		return NewMysqlBackupService()
	}
	return nil, errors.New("Unspecified database type = " + database_type)
}
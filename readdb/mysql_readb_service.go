package readdb

import (
	"os"
	"errors"
 	"database/sql"
 	_"github.com/go-sql-driver/mysql"
	"database-backup-manager/constants"
)

type mysqlReaddbService struct {
	db *sql.DB
}

var exclude_dbs = [...]string{"information_schema", "performance_schema", "mysql", "sys"}

func NewMysqlReaddbService() (*mysqlReaddbService,error) {
	host := os.Getenv(constants.DATABASE_HOST)
	if (len(host) <= 0) {return nil, errors.New("No database host specified!")}
	port := os.Getenv(constants.DATABASE_PORT)
	if (len(port) <= 0) {port = "3306"}
	user := os.Getenv(constants.DATABASE_USER)
	if (len(user) <= 0) {return nil, errors.New("No database user specified!")}
	password := os.Getenv(constants.DATABASE_PASSWORD)
	if (len(password) <= 0) {return nil, errors.New("No database password specified!")}

	db, err := sql.Open("mysql", user + ":" + password + "@tcp(" + host + ":" + port + ")/")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
 	   return nil, err
	}

	drs := new(mysqlReaddbService)
	drs.db = db
	return drs,nil
}

func (dr *mysqlReaddbService) ReadDatabaseNames() ([]string, error){
	stmtOut, err := dr.db.Prepare("SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer stmtOut.Close()

	results, err := stmtOut.Query();
	if err != nil {
		return nil, err
	}

	databaseNames := []string{}
	for results.Next() {
		var dbName = "";
		results.Scan(&dbName);
		if (dr.isValidDb(dbName)) {
			databaseNames = append(databaseNames, dbName)
		}
	}
	return databaseNames, nil
}
/*
func (dr *mysqlReaddbService) GetDatabases() ([]string, error) {
	db, err := dr.connect()
	if (err != nil) {
		return nil, err
	}
	defer db.Close()

	return dr.readDatabaseNames(db)
}*/

func (dr *mysqlReaddbService) isValidDb(dbname string) bool {
    for _, excludeDb := range exclude_dbs {
        if excludeDb == dbname {
            return false
        }
    }
    return true
}
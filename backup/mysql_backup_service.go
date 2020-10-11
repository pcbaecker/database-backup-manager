package backup

import (
	"fmt"
	"os"
	"bytes"
	"io/ioutil"
	"os/exec"
	"archive/zip"
	"database-backup-manager/constants"
)

type mysqlBackupService struct {
	host string
	port string
	user string
	password string
}

func NewMysqlBackupService() (*mysqlBackupService,error) {
	dbs := new(mysqlBackupService)
	dbs.host = os.Getenv(constants.DATABASE_HOST)
	dbs.port = os.Getenv(constants.DATABASE_PORT)
	if (len(dbs.port) <= 0) {dbs.port = "3306"}
	dbs.user = os.Getenv(constants.DATABASE_USER)
	dbs.password = os.Getenv(constants.DATABASE_PASSWORD)
	return dbs,nil
}

func (dbs *mysqlBackupService) BackupDatabases(databaseNames []string) (string,error) {
	dir, err := ioutil.TempDir(".", "databaseBackup*")
	if (err != nil) {return "", err}
	defer os.RemoveAll(dir)

	zipFile, err := ioutil.TempFile(".", "tmpzip")
    if err != nil {return "", err}
    defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, dbname := range databaseNames {
		filename := dbname + ".sql"
		filepath := dir + "/" + filename
		err = dbs.backupDatabase(filepath, dbname)
		if (err == nil) {
			dbs.addFileToZip(zipWriter, filepath, filename)
		}
	}
	
	return zipFile.Name(), nil
}

func (dbs *mysqlBackupService) addFileToZip(zipWriter *zip.Writer, filepath string, filename string) error {
	zipFileHandle, err := zipWriter.Create(filename)
	if (err != nil) {return err}
	filebody, err2 := ioutil.ReadFile(filepath)
	if (err2 != nil) {return err2}
	_, err = zipFileHandle.Write([]byte(filebody))
	if err != nil {return err}

	return nil
}

func (dbs *mysqlBackupService) backupDatabase(filename string, databaseName string) (error) {
	cmd := exec.Command("mysqldump", "--set-gtid-purged=OFF", "--column-statistics=0", "--host=" + dbs.host, "--port=" + dbs.port, "--user=" + dbs.user, "--password=" + dbs.password, "--databases", databaseName)
	
	outfile, err := os.Create(filename)
    if err != nil {
        return err
    }
	defer outfile.Close()
	var errb bytes.Buffer
	cmd.Stdout = outfile
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		fmt.Printf("(db=%s) cmd.Run() failed with %s\n%s\n", databaseName, err, errb.String())
		return err
	}

	return nil
}
package main

import (
    "os"
    "fmt"
    "time"
    "strconv"
    "database-backup-manager/backup"
    "database-backup-manager/readdb"
    "database-backup-manager/storage"
    "database-backup-manager/notification"
)

type result struct {
    numberOfDatabases int
    targetFilename string
    warnings bool
}

func execute(
    readdbService readdb.ReaddbService,
    databaseBackupService backup.DatabaseBackupService,
    storageService storage.StorageService) (*result, error) {
        result := &result{}

        fmt.Println("Reading database names")
        databaseNames, err := readdbService.ReadDatabaseNames();
        if (err != nil) {return result, nil}
        result.numberOfDatabases = len(databaseNames)
        fmt.Println(strconv.Itoa(result.numberOfDatabases) + " databases found!")
    
        fmt.Println("Starting backup of databases")
        zipFile, warnings, err := databaseBackupService.BackupDatabases(databaseNames)
        if (err != nil) {return result, nil}
        defer os.Remove(zipFile)
        result.warnings = warnings
        fmt.Println("Databases backuped and zipped")
    
        fmt.Println("Starting upload of backupfile")
        result.targetFilename = "backup-" + time.Now().Format(time.RFC3339) + ".zip"
        err = storageService.UploadFile(zipFile, result.targetFilename, "application/zip")
        if (err != nil) {return result, nil}
        fmt.Println("Backupfile successfully uploaded with name = " + result.targetFilename)

        return result, nil
}

func main() {
    fmt.Println("Preparing services")
    notificationService, err := notification.NewNotificationService()
    if (err != nil) {
        panic(err.Error())
    }
    readdbService, err := readdb.NewReaddbService()
    if (err != nil) {
        panic(err.Error())
    }
    databaseBackupService, err := backup.NewDatabaseBackupService()
    if (err != nil) {
        panic(err.Error())
    }
    storageService, err := storage.NewStorageService()
    if (err != nil) {
        panic(err.Error())
    }
    fmt.Println("All services loaded")

    start := time.Now()
    result, err := execute(readdbService, databaseBackupService, storageService)
    duration := time.Since(start)

    if (err != nil) {
        msg := "Error: Could not create backup: " + err.Error()
        fmt.Println(msg)
        notificationService.Send(msg)
    } else {
        msg := "Successfully uploaded backupfile " + result.targetFilename + " with " + strconv.Itoa(result.numberOfDatabases) + " databases! (" + duration.String() + " )"
        fmt.Println(msg)
        notificationService.Send(msg)
    }
}

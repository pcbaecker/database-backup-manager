# Database Backup Manager

This database backup manager can be used to backup multiple databases to an object storage. It is designed to be run regulary by a cronjob or scheduled job with kubernetes. The configuration is completly done with environment variables.

## Tips

The mysql part is tested with mariadb. That may be important due to the difference in the mysqldump implementation.

    sudo apt install mariadb-client

## Environment variables

    DATABASE_TYPE = "mysql"
    DATABASE_HOST = "1.2.3.4"
    DATABASE_PORT = "3306" (optional)
    DATABASE_USER = "root"
    DATABASE_PASSWORD = "myrootpassword"

    STORAGE_TYPE = "gcs"
    STORAGE_SECRETKEY = "e0FTMiJRMjNRMjMi..." (base64 of the json credentials)
    STORAGE_BUCKET = "mybackupbucket"

    NOTIFICATION_TYPE = "slack" (optional)
    SLACK_WEBHOOK_URL = "https://hooks.slack.com/services/T01BXTS7JFM/B01CJV4M9UL/" (optional)

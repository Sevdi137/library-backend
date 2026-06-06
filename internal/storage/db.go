package storage

import (
    "database/sql"
    "log/slog"
    "os"

    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) error {
    var err error
    DB, err = sql.Open("sqlite3", dbPath)
    if err != nil {
        return err
    }

    if err = DB.Ping(); err != nil {
        return err
    }

    slog.Info("database connected", "path", dbPath)
    return nil
}

func CloseDB() {
    if DB != nil {
        DB.Close()
    }
}

func RunMigrations(migrationPath string) error {
    migrationSQL, err := os.ReadFile(migrationPath)
    if err != nil {
        return err
    }

    _, err = DB.Exec(string(migrationSQL))
    return err
}

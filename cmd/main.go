package main

import (
    "log/slog"
    "net/http"
    "os"

    "github.com/joho/godotenv"

    "library-backend/internal/auth"
    "library-backend/internal/handlers"
    "library-backend/internal/middleware"
    "library-backend/internal/storage"
)

func main() {
    if err := godotenv.Load(); err != nil {
        slog.Warn("no .env file found, using environment variables")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    dbPath := os.Getenv("DB_PATH")
    if dbPath == "" {
        dbPath = "./library.db"
    }
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        slog.Warn("JWT_SECRET not set, using default (INSECURE for production)")
        jwtSecret = "default-insecure-secret"
    }

    if err := storage.InitDB(dbPath); err != nil {
        slog.Error("failed to connect to database", "error", err)
        os.Exit(1)
    }
    defer storage.CloseDB()

    if err := storage.RunMigrations("migrations/001_init.sql"); err != nil {
        slog.Error("failed to run migrations", "error", err)
        os.Exit(1)
    }

    auth.SetSecret(jwtSecret)

    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)

    http.HandleFunc("POST /login", handlers.Login)
    http.HandleFunc("GET /books", handlers.ListBooks)
    http.HandleFunc("GET /books/{id}", handlers.GetBook)

    http.HandleFunc("POST /books", middleware.Auth(handlers.CreateBook, "admin"))
    http.HandleFunc("PUT /books/{id}", middleware.Auth(handlers.UpdateBook, "admin"))
    http.HandleFunc("POST /users", middleware.Auth(handlers.RegisterUser, "admin"))

    http.HandleFunc("GET /users/{id}/books", middleware.Auth(handlers.GetUserBooks, ""))
    http.HandleFunc("POST /issues", middleware.Auth(handlers.IssueBook, ""))
    http.HandleFunc("POST /returns", middleware.Auth(handlers.ReturnBook, ""))

    slog.Info("server starting", "port", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        slog.Error("server failed", "error", err)
        os.Exit(1)
    }
}

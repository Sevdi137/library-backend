package handlers

import (
    "encoding/json"
    "net/http"
    "strings"
    "time"

    "library-backend/internal/auth"
    "library-backend/internal/models"
    "library-backend/internal/storage"

    "github.com/google/uuid"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
        return
    }

    if user.Name == "" || user.Email == "" || user.Password == "" {
        http.Error(w, `{"error":"name, email and password are required"}`, http.StatusBadRequest)
        return
    }

    user.ID = uuid.New().String()
    user.RegistrationDate = time.Now()
    if user.Role == "" {
        user.Role = "user"
    }

    if err := storage.CreateUser(user); err != nil {
        if strings.Contains(err.Error(), "UNIQUE") {
            http.Error(w, `{"error":"user with this email already exists"}`, http.StatusConflict)
            return
        }
        http.Error(w, `{"error":"failed to create user"}`, http.StatusInternalServerError)
        return
    }

    user.Password = ""
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func GetUserBooks(w http.ResponseWriter, r *http.Request) {
    userID := strings.TrimPrefix(r.URL.Path, "/users/")
    userID = strings.TrimSuffix(userID, "/books")
    
    if userID == "" {
        http.Error(w, `{"error":"user id required"}`, http.StatusBadRequest)
        return
    }

    if err := storage.ValidateUserExists(userID); err != nil {
        http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
        return
    }

    books, err := storage.GetUserBooks(userID)
    if err != nil {
        http.Error(w, `{"error":"failed to fetch user books"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "user_id": userID,
        "books":   books,
        "count":   len(books),
    })
}

func Login(w http.ResponseWriter, r *http.Request) {
    var creds models.Credentials
    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
        return
    }

    user, err := storage.GetUserByEmail(creds.Email)
    if err != nil {
        http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
        return
    }

    if user.Password != creds.Password {
        http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
        return
    }

    token, err := auth.GenerateToken(user.ID, user.Role)
    if err != nil {
        http.Error(w, `{"error":"failed to generate token"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "token": token,
        "role":  user.Role,
    })
}

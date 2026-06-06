package middleware

import (
    "context"
    "net/http"
    "strings"
    "library-backend/internal/auth"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "userRole"

func Auth(next http.HandlerFunc, requiredRole string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "missing token", http.StatusUnauthorized)
            return
        }
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "invalid token format", http.StatusUnauthorized)
            return
        }
        userID, role, err := auth.ValidateToken(parts[1])
        if err != nil {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }
        if requiredRole != "" && role != requiredRole {
            http.Error(w, "forbidden", http.StatusForbidden)
            return
        }
        ctx := context.WithValue(r.Context(), UserIDKey, userID)
        ctx = context.WithValue(ctx, UserRoleKey, role)
        next(w, r.WithContext(ctx))
    }
}

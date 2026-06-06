package auth

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func SetSecret(secret string) {
    jwtSecret = []byte(secret)
}

func GenerateToken(userID, role string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "role":    role,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
        "iat":     time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (string, string, error) {
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return jwtSecret, nil
    })

    if err != nil {
        return "", "", err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return "", "", errors.New("invalid token")
    }

    userID, ok := claims["user_id"].(string)
    if !ok {
        return "", "", errors.New("user_id not found")
    }

    role, ok := claims["role"].(string)
    if !ok {
        role = "user"
    }

    return userID, role, nil
}

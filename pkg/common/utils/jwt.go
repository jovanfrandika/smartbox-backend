package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jovanfrandika/smartbox-backend/pkg/common/config"
)

func CreateToken(ttl time.Duration, payload interface{}, secretKey string) (string, error) {
	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))

	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func ValidateToken(token string, secretKey string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return "", fmt.Errorf("validate: invalid token")
	}

	return fmt.Sprintf("%s", claims["sub"]), nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		}
		jwtToken := authHeader[1]
		userId, err := ValidateToken(jwtToken, config.Cfg.JWTAccessSecretKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "userId", userId)))
	})
}

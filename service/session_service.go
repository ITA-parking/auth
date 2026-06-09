package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strings"
	"time"
)

type ctxKey string

var userIDCtxKey ctxKey = "userId"

var jwtSecret = []byte(os.Getenv("SERVICE_JWT_SECRET"))

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uuid.UUID, expiresIn time.Duration) (string, error) {
	userIdString := userID.String()
	claims := Claims{
		UserID: userIdString,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   userIdString,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyJWT(jwtToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func AuthMiddleware(handler http.Handler) http.Handler {
	excludedPaths := []string{"/auth", "/docs", "/openapi.yaml", "/openapi.json"}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for excluded paths
		for _, path := range excludedPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				handler.ServeHTTP(w, r)
				return
			}
		}

		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := VerifyJWT(tokenStr)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		userId, idParseErr := uuid.Parse(claims.UserID)
		if idParseErr != nil {
			http.Error(w, "invalid or expired user id", http.StatusUnauthorized)
		}
		// Add userID to request context
		ctx := context.WithValue(r.Context(), userIDCtxKey, userId)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) *uuid.UUID {
	if v, ok := ctx.Value(userIDCtxKey).(uuid.UUID); ok {
		return &v
	}

	return nil
}

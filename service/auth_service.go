package service

import (
	"auth-service/dto"
	"auth-service/repo"
	"auth-service/repo/model"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	argonTime    = 10        // iterations
	argonMemory  = 64 * 1024 // 64 MB
	argonThreads = 4
	argonKeyLen  = 32
	saltLen      = 16

	refreshTokenExpiry = 30 * 24 * time.Hour // 30 days
	accessTokenExpiry  = time.Hour
)

func LoginUser(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	if req == nil {
		return nil, huma.Error400BadRequest("invalid request")
	}

	if req.Body.Username == "" || req.Body.Password == "" {
		return nil, huma.Error400BadRequest("username or password is empty")
	}

	var dbUser *model.User
	var err error
	if strings.Contains(req.Body.Username, "@") {
		dbUser, err = repo.FindUserByEmail(context.Background(), req.Body.Username)
	} else {
		dbUser, err = repo.FindUserByUsername(context.Background(), req.Body.Username)
	}

	if err != nil || dbUser == nil {
		return nil, huma.Error404NotFound("user not found", err)
	}

	if err = VerifyPassword(dbUser.Password, req.Body.Password); err != nil {
		return nil, huma.Error401Unauthorized("invalid password or username")
	}

	jwt, jwtErr := GenerateJWT(dbUser.ID, accessTokenExpiry)
	if jwtErr != nil {
		slog.Error("Error generating JWT token", "error", jwtErr.Error())
		return nil, huma.Error500InternalServerError("Internal Server Error", jwtErr)
	}

	refreshToken, rtErr := generateRefreshToken(context.Background(), dbUser.ID)
	if rtErr != nil {
		return nil, huma.Error500InternalServerError("Internal Server Error", rtErr)
	}

	return &dto.LoginResponse{
		Status: http.StatusOK,
		Body: dto.LoginResponseBody{
			AccessToken:  jwt,
			TokenType:    "Bearer",
			ExpiresIn:    int(accessTokenExpiry.Seconds()),
			RefreshToken: refreshToken,
		},
	}, nil
}

func RegisterNewUser(req *dto.RegisterRequest) error {
	validationErr := validateRegisterRequest(req)
	if validationErr != nil {
		return huma.Error400BadRequest("Request validation error", validationErr)
	}

	passwordHash, hashErr := HashPassword(req.Body.Password)
	if hashErr != nil {
		return huma.Error500InternalServerError("Internal Server Error", hashErr)
	}

	userID := uuid.New()
	user := model.User{
		ID:       userID,
		Username: req.Body.Username,
		Email:    req.Body.Email,
		Password: passwordHash,
	}

	if saveErr := repo.SaveUser(context.Background(), user); saveErr != nil {
		return huma.Error500InternalServerError("Internal Server Error", saveErr)
	}

	return nil
}

func RefreshSession(req *dto.RefreshRequest) (*dto.LoginResponse, error) {
	if req == nil || req.Body.RefreshToken == "" {
		return nil, huma.Error400BadRequest("refresh_token is required")
	}

	stored, err := repo.FindRefreshToken(context.Background(), req.Body.RefreshToken)
	if err != nil {
		return nil, huma.Error500InternalServerError("Internal Server Error", err)
	}
	if stored == nil || stored.Revoked || time.Now().After(stored.ExpiresAt) {
		return nil, huma.Error401Unauthorized("invalid or expired refresh token")
	}

	if err = repo.RevokeRefreshToken(context.Background(), req.Body.RefreshToken); err != nil {
		return nil, huma.Error500InternalServerError("Internal Server Error", err)
	}

	jwt, jwtErr := GenerateJWT(stored.UserID, accessTokenExpiry)
	if jwtErr != nil {
		return nil, huma.Error500InternalServerError("Internal Server Error", jwtErr)
	}

	newRefreshToken, rtErr := generateRefreshToken(context.Background(), stored.UserID)
	if rtErr != nil {
		return nil, huma.Error500InternalServerError("Internal Server Error", rtErr)
	}

	return &dto.LoginResponse{
		Status: http.StatusOK,
		Body: dto.LoginResponseBody{
			AccessToken:  jwt,
			TokenType:    "Bearer",
			ExpiresIn:    int(accessTokenExpiry.Seconds()),
			RefreshToken: newRefreshToken,
		},
	}, nil
}

func generateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	tokenStr := hex.EncodeToString(b)

	rt := model.RefreshToken{
		ID:        uuid.New(),
		Token:     tokenStr,
		UserID:    userID,
		ExpiresAt: time.Now().Add(refreshTokenExpiry),
		Revoked:   false,
	}

	if err := repo.SaveRefreshToken(ctx, rt); err != nil {
		return "", err
	}

	return tokenStr, nil
}

func validateRegisterRequest(req *dto.RegisterRequest) error {
	if req.Body.Username == "" {
		return errors.New("username is empty")
	}

	if req.Body.Password == "" {
		return errors.New("password is empty")
	}

	if req.Body.Email == "" {
		return errors.New("email is empty")
	}

	if strings.Contains(req.Body.Username, "@") {
		return errors.New("username must not contain @")
	}

	if strings.Contains(req.Body.Username, " ") {
		return errors.New("username must not contain space")
	}

	if !strings.Contains(req.Body.Email, "@") {
		return errors.New("invalid email")
	}

	return nil
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory,
		argonTime,
		argonThreads,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

func VerifyPassword(encoded, password string) error {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return errors.New("invalid hash format")
	}

	var memory uint32
	var iterations uint32
	var threads uint8

	_, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&memory,
		&iterations,
		&threads,
	)
	if err != nil {
		return err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return err
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		threads,
		uint32(len(expectedHash)),
	)

	if subtle.ConstantTimeCompare(hash, expectedHash) == 1 {
		return nil
	}

	return errors.New("invalid password")
}

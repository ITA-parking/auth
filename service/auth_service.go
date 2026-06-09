package service

import (
	"auth-service/dto"
	"auth-service/repo"
	"auth-service/repo/model"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
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
		//this could technically also be an 500 internal server error due to db error
		return nil, huma.Error404NotFound("user not found", err)
	}

	passwordValidationErr := VerifyPassword(dbUser.Password, req.Body.Password)
	if passwordValidationErr != nil {
		return nil, huma.Error401Unauthorized("invalid password or username")
	}

	//FIXME: actually create jwt token and refresh token!!

	expiresIn := time.Minute * 60
	jwt, jwtErr := GenerateJWT(dbUser.ID, expiresIn)
	if jwtErr != nil {
		slog.Error("Error generating JWT token", "error", jwtErr.Error())
		return nil, huma.Error500InternalServerError("Internal Server Error", jwtErr)
	}

	return &dto.LoginResponse{
		Status: http.StatusOK,
		Body: dto.LoginResponseBody{
			AccessToken:  jwt,
			TokenType:    "Bearer",
			ExpiresIn:    int(expiresIn.Seconds()),
			RefreshToken: "REFRESH_TOKEN", //TODO!!!
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
	user := model.User{
		ID:       uuid.New(),
		Username: req.Body.Username,
		Email:    req.Body.Email,
		Password: passwordHash,
	}

	saveErr := repo.SaveUser(context.Background(), user)
	if saveErr != nil {
		return huma.Error500InternalServerError("Internal Server Error", saveErr)
	}

	return nil
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
		return errors.New("invalida email")
	}

	if req.Body.EncPrivKey == nil || len(req.Body.EncPrivKey) == 0 {
		return errors.New("encrypt private key is empty")
	}

	if req.Body.PubKey == nil || len(req.Body.PubKey) == 0 {
		return errors.New("public key is empty")
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
	var time uint32
	var threads uint8

	_, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&memory,
		&time,
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
		time,
		memory,
		threads,
		uint32(len(expectedHash)),
	)

	if subtle.ConstantTimeCompare(hash, expectedHash) == 1 {
		return nil
	}

	return errors.New("invalid password")
}

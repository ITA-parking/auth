package service

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

func init() {
	os.Setenv("SERVICE_JWT_SECRET", "test-secret-key-for-unit-tests")
	jwtSecret = []byte(os.Getenv("SERVICE_JWT_SECRET"))
}

func TestGenerateAndVerifyJWT(t *testing.T) {
	userID := uuid.New()

	token, err := GenerateJWT(userID, time.Hour)
	if err != nil {
		t.Fatalf("GenerateJWT returned error: %s", err)
	}
	if token == "" {
		t.Fatal("GenerateJWT returned empty token")
	}

	claims, err := VerifyJWT(token)
	if err != nil {
		t.Fatalf("VerifyJWT returned error: %s", err)
	}
	if claims.UserID != userID.String() {
		t.Errorf("expected UserID %s, got %s", userID.String(), claims.UserID)
	}
}

func TestVerifyJWT_Expired(t *testing.T) {
	userID := uuid.New()

	token, err := GenerateJWT(userID, -time.Second)
	if err != nil {
		t.Fatalf("GenerateJWT returned error: %s", err)
	}

	_, err = VerifyJWT(token)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestVerifyJWT_InvalidToken(t *testing.T) {
	_, err := VerifyJWT("this.is.not.a.valid.jwt")
	if err == nil {
		t.Error("expected error for invalid token string, got nil")
	}
}

func TestVerifyJWT_EmptyToken(t *testing.T) {
	_, err := VerifyJWT("")
	if err == nil {
		t.Error("expected error for empty token, got nil")
	}
}

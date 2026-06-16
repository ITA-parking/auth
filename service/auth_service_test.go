package service

import (
	"auth-service/dto"
	"testing"
)

func TestHash(t *testing.T) {
	plainText := "TESTPASSWORD"
	hash, hashErr := HashPassword(plainText)
	if hashErr != nil {
		t.Errorf("Error hashing password: %s", hashErr)
	}

	validPasswordResult := VerifyPassword(hash, plainText)
	if validPasswordResult != nil {
		t.Errorf("Error verifying password: %s", validPasswordResult)
	}

	invalidPasswordResult := VerifyPassword(hash, "invalidpass")
	if invalidPasswordResult == nil {
		t.Errorf("Error invalid password verified as correct: %s", invalidPasswordResult)
	}
}

func TestValidateRegisterRequest_Valid(t *testing.T) {
	req := &dto.RegisterRequest{}
	req.Body.Username = "alice"
	req.Body.Password = "secret123"
	req.Body.Email = "alice@example.com"

	if err := validateRegisterRequest(req); err != nil {
		t.Errorf("expected nil error for valid request, got: %s", err)
	}
}

func TestValidateRegisterRequest_MissingUsername(t *testing.T) {
	req := &dto.RegisterRequest{}
	req.Body.Username = ""
	req.Body.Password = "secret"
	req.Body.Email = "a@b.com"

	if err := validateRegisterRequest(req); err == nil {
		t.Error("expected error for empty username")
	}
}

func TestValidateRegisterRequest_MissingPassword(t *testing.T) {
	req := &dto.RegisterRequest{}
	req.Body.Username = "alice"
	req.Body.Password = ""
	req.Body.Email = "a@b.com"

	if err := validateRegisterRequest(req); err == nil {
		t.Error("expected error for empty password")
	}
}

func TestValidateRegisterRequest_MissingEmail(t *testing.T) {
	req := &dto.RegisterRequest{}
	req.Body.Username = "alice"
	req.Body.Password = "secret"
	req.Body.Email = ""

	if err := validateRegisterRequest(req); err == nil {
		t.Error("expected error for empty email")
	}
}

func TestValidateRegisterRequest_UsernameWithAt(t *testing.T) {
	req := &dto.RegisterRequest{}
	req.Body.Username = "ali@ce"
	req.Body.Password = "secret"
	req.Body.Email = "a@b.com"

	if err := validateRegisterRequest(req); err == nil {
		t.Error("expected error for username containing @")
	}
}

func TestValidateRegisterRequest_UsernameWithSpace(t *testing.T) {
	req := &dto.RegisterRequest{}
	req.Body.Username = "ali ce"
	req.Body.Password = "secret"
	req.Body.Email = "a@b.com"

	if err := validateRegisterRequest(req); err == nil {
		t.Error("expected error for username containing space")
	}
}

func TestValidateRegisterRequest_InvalidEmail(t *testing.T) {
	req := &dto.RegisterRequest{}
	req.Body.Username = "alice"
	req.Body.Password = "secret"
	req.Body.Email = "notanemail"

	if err := validateRegisterRequest(req); err == nil {
		t.Error("expected error for email without @")
	}
}

package service

import "testing"

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

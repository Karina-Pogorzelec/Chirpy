package auth

import (
	"testing"
	"time"
	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "my-secret-password"
    
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 3. Ensure the hash isn't empty and isn't the same as the password
	if hash == "" {
		t.Fatalf("expected non-empty hash")
	}
	if hash == password {
		t.Fatalf("expected hash to be different from password")
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("expected no error during check, got %v", err)
	}
	if !match {
		t.Errorf("expected password to match the hash, but it didn't")
	}

	wrongPassword := "not-my-password"
	match, err = CheckPasswordHash(wrongPassword, hash)
	if err != nil {
		t.Fatalf("expected no error during check of wrong password, got %v", err)
	}
	if match {
		t.Errorf("expected wrong password NOT to match, but it did")
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my-secret-key"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	returnedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if returnedUserID != userID {
		t.Errorf("expected userID %v, got %v", userID, returnedUserID)
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
    userID := uuid.New()
    tokenSecret := "correct-secret"
    wrongSecret := "wrong-secret"
    expiresIn := time.Hour

    // 1. Create a token using the 'tokenSecret'
    token, _ := MakeJWT(userID, tokenSecret, expiresIn)

    // 2. Try to validate it using the 'wrongSecret'
    _, err := ValidateJWT(token, wrongSecret)

    // 3. Check if 'err' is NOT nil
    if err == nil {
        t.Errorf("expected error when validating with wrong secret, got nil")
    }
}

func TestValidateJWT_Expired(t *testing.T) {
    userID := uuid.New()
    tokenSecret := "my-secret"
    // Set an expiration that has already passed
    expiresIn := -time.Hour 

    token, _ := MakeJWT(userID, tokenSecret, expiresIn)

    _, err := ValidateJWT(token, tokenSecret)

    // How should we check 'err' here to confirm it failed specifically because of expiration?
    if err == nil {
        t.Errorf("expected error when validating expired token, got nil")
    }
}
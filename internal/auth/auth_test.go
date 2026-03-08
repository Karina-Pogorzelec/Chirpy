package auth

import (
	"testing"
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
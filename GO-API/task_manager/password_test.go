package main

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("123456")
	if err != nil {
		t.Fatal(err)
	}

	if hash == "" {
		t.Fatal("hash vazio")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("123456")); err != nil {
		t.Fatalf("hash não corresponde à senha original: %v", err)
	}
}

package password

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")

	if err != nil {
		t.Errorf("expected nil err, found %v", err)
	}

	if len(hash) == 0 {
		t.Errorf("hash has len 0, expected greater than or equal to 1")
	}
}

func TestCheckPassword(t *testing.T) {
	hash, err := HashPassword("password")

	if err != nil {
		t.Errorf("expected nil err, found %v", err)
	}

	if !CheckPassword("password", hash) {
		t.Error("expected 'password' to match hash")
	}

	if CheckPassword("wrong", hash) {
		t.Error("expected 'wrong' to not match hash")
	}
}

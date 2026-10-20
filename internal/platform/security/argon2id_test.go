package security

import (
	"testing"
)

func TestArgon2Hasher(t *testing.T) {
	// Use small parameters for fast test execution
	fastParams := Argon2Params{
		Memory:      1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}
	hasher := NewArgon2Hasher(fastParams)

	pw := "SecretPassword123!"
	encoded, err := hasher.Hash(pw)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	match, needsRehash, err := hasher.Verify(pw, encoded)
	if err != nil {
		t.Fatalf("verify failed with error: %v", err)
	}
	if !match {
		t.Errorf("expected password to match")
	}
	if needsRehash {
		t.Errorf("expected needsRehash=false")
	}

	// Wrong password
	match, _, err = hasher.Verify("WrongPassword!", encoded)
	if err != nil {
		t.Fatalf("verify failed with error: %v", err)
	}
	if match {
		t.Errorf("expected wrong password not to match")
	}

	// Parameter upgrade detection
	upgradedHasher := NewArgon2Hasher(Argon2Params{
		Memory:      2048,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})
	match, needsRehash, err = upgradedHasher.Verify(pw, encoded)
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if !match {
		t.Errorf("expected match")
	}
	if !needsRehash {
		t.Errorf("expected needsRehash=true due to memory change")
	}
}


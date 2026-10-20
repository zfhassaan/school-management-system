package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultArgon2Params gives RFC 9106 recommended interactive parameters (64 MB, 3 iterations, 2 threads, 32-byte key)
var DefaultArgon2Params = Argon2Params{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

type Argon2Hasher struct {
	params Argon2Params
}

func NewArgon2Hasher(params Argon2Params) *Argon2Hasher {
	return &Argon2Hasher{params: params}
}

func NewDefaultArgon2Hasher() *Argon2Hasher {
	return &Argon2Hasher{params: DefaultArgon2Params}
}

// Hash creates a PHC string format hash: $argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
func (h *Argon2Hasher) Hash(password string) (string, error) {
	salt := make([]byte, h.params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.params.Iterations,
		h.params.Memory,
		h.params.Parallelism,
		h.params.KeyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.params.Memory,
		h.params.Iterations,
		h.params.Parallelism,
		b64Salt,
		b64Hash,
	)
	return encoded, nil
}

// Verify compares a candidate password against the PHC encoded hash.
func (h *Argon2Hasher) Verify(password, encodedHash string) (match bool, needsRehash bool, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return false, false, errors.New("invalid hash format")
	}

	if vals[1] != "argon2id" {
		return false, false, errors.New("unsupported hash algorithm")
	}

	var version int
	if _, err := fmt.Sscanf(vals[2], "v=%d", &version); err != nil || version != argon2.Version {
		return false, false, errors.New("incompatible argon2 version")
	}

	var memory, iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return false, false, errors.New("invalid hash parameters")
	}

	salt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return false, false, fmt.Errorf("invalid salt encoding: %w", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return false, false, fmt.Errorf("invalid hash encoding: %w", err)
	}

	keyLen := uint32(len(hash))
	comparisonHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLen)

	if subtle.ConstantTimeCompare(hash, comparisonHash) != 1 {
		return false, false, nil
	}

	needsRehash = memory != h.params.Memory ||
		iterations != h.params.Iterations ||
		parallelism != h.params.Parallelism ||
		keyLen != h.params.KeyLength

	return true, needsRehash, nil
}

// GenerateRandomToken generates a cryptographically secure URL-safe base64 string
func GenerateRandomToken(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}


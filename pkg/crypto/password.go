package crypto

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return password, nil
	}

	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	hash = append(salt, hash...)

	return hex.EncodeToString(hash[:]), nil
}

func CheckPassword(password, storedHash string) (bool, error) {
	storedBytes, err := hex.DecodeString(storedHash)
	if err != nil {
		return false, err
	}

	salt := storedBytes[0:8]
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	hash = append(salt, hash...)

	return bytes.Equal(hash, storedBytes), nil
}

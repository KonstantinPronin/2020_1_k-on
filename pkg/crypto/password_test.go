package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckPassword(t *testing.T) {
	password := "test"

	hash, err := HashPassword(password)
	assert.Nil(t, err)

	res, err := CheckPassword(password, hash)
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestCheckPassword_WrongPassword(t *testing.T) {
	password := "test"
	wrongPassword := "test1"

	hash, err := HashPassword(password)
	assert.Nil(t, err)

	res, err := CheckPassword(wrongPassword, hash)
	assert.Nil(t, err)
	assert.False(t, res)
}

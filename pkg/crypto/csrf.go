package crypto

import (
	"encoding/hex"
	"golang.org/x/crypto/argon2"
)

const (
	CSRFHeader = "X-CSRF-TOKEN"
	Key        = "eE%yh?aAH_hYk*5h$DXvTddAGt2eWCt^+TT_4*$ADxz^X$5ue74jmeJT@z^+c_*v"
)

func CreateToken(sessionId string) string {
	hash := argon2.IDKey([]byte(Key), []byte(sessionId), 1, 64*1024, 4, 32)
	return hex.EncodeToString(hash[:])
}

func CheckToken(sessionId, token string) bool {
	return token == CreateToken(sessionId)
}

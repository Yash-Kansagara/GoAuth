package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func GetHash(val string) string {
	salt := make([]byte, 16)
	rand.Read(salt)
	hash := argon2.IDKey([]byte(val), salt, 1, 64*1024, 4, 32)
	saltb64 := base64.StdEncoding.EncodeToString(salt)
	hashb64 := base64.StdEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s:%s", saltb64, hashb64)
}

func GetHashWithSalt(saltString string, val string) string {
	salt, err := base64.StdEncoding.DecodeString(saltString)
	if err != nil {
		return ""
	}
	hash := argon2.IDKey([]byte(val), salt, 1, 64*1024, 4, 32)
	return base64.StdEncoding.EncodeToString(hash)
}

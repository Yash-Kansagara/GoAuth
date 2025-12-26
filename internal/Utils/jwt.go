package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignToken(username string, userid string, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	expiray := os.Getenv("JWT_EXPIRY")

	claims := jwt.MapClaims{
		"userid":   userid,
		"username": username,
		"email":    email,
	}

	if jwtExpiresin, err := time.ParseDuration(expiray); err != nil {
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	} else {
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(jwtExpiresin))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTokem, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("Error signing token. %s", err.Error())
	}

	return signedTokem, nil
}

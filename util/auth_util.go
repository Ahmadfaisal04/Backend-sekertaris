package util

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(email, role string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey, err := generateRandomKey(32)
	if err != nil {
		return "", err
	}

	return token.SignedString(jwtKey)
}

func generateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	return key, err
}
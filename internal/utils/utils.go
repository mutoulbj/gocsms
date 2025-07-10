package utils

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func ParseUUID(idStr string) (uuid.UUID, error) {
	// This function is a placeholder for UUID parsing logic.
	// In a real implementation, you would use a library like "github.com/google/uuid"
	// to parse and validate the UUID.
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

// HashPassword hashes the password with a salt using bcrypt
func HashPassword(password, salt string) (string, error) {
	combined := password + salt

	// use bcrypt to hash the password with the salt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// VerifyPassword verifies if the provided password is correct
func VerifyPassword(password, salt, hashedPassword string) bool {
	combined := password + salt

	// use bcrypt to verify the password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(combined))
	return err == nil
}

// GenerateSalt to create a random salt for password hashing
func GenerateSalt() string {
	return GenerateRandomString(16)
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	// calculate the byte length needed for base64 encoding
	byteLength := max((length*3)/4, length)

	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		// if failed, use charset method
		return randomStringWithCharset(length)
	}

	// use base64 RawURLEncoding to avoid padding and URL-safe characters
	encoded := base64.RawURLEncoding.EncodeToString(bytes)
	if len(encoded) > length {
		return encoded[:length]
	}
	return encoded
}

func randomStringWithCharset(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}
	return string(result)
}

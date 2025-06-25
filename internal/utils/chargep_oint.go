package utils

import (
	"math/rand"
)

func GenerateSerialNumber() string {
	// Generate a unique serial number for the charge point
	return "SN-" + randomString(10)
}

func randomString(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

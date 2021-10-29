package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateValidationCode(length int) string {
	code := ""
	for i := 0; i < length; i++ {
		result, _ := rand.Int(rand.Reader, big.NewInt(100))
		code += string(rune(result.Int64()%10 + 48))
	}

	return code
}

package utils

import (
	"math/rand"
	"time"
)

func GenerateValidationCode(length int) string {
	numbers := "0123456789"

	rand.Seed(time.Now().UnixNano())

	var code string

	for i := 0; i < length; i++ {
		code += string(numbers[rand.Intn(len(numbers))])
	}

	return code
}

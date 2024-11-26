package services

import "strings"

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Encode(number int64) string {
	if number == 0 {
		return string(alphabet[0])
	}

	var encoded strings.Builder
	base := int64(len(alphabet))

	for number > 0 {
		remainder := number % base
		encoded.WriteByte(alphabet[remainder])
		number = number / base
	}
	// Reversing the string
	result := encoded.String()
	runes := []rune(result)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

package services

import (
	"math/big"
	"strings"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// EncodeBigInt encodes a big.Int number into a base62 string
func EncodeBigInt(number *big.Int) string {
	if number.Cmp(big.NewInt(0)) == 0 {
		return string(alphabet[0])
	}

	var encoded strings.Builder
	base := big.NewInt(int64(len(alphabet)))
	zero := big.NewInt(0)

	for number.Cmp(zero) > 0 {
		remainder := new(big.Int)
		number.DivMod(number, base, remainder)
		encoded.WriteByte(alphabet[remainder.Int64()])
	}

	// Reverse the string
	result := encoded.String()
	runes := []rune(result)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

package services

import (
	"errors"
	"math/big"
	"strings"
)

// EncodeBigInt encodes a big.Int to a base62 string.
// It returns an error if the input is nil or negative.
func EncodeBigInt(n *big.Int) (string, error) {
	if n == nil {
		return "", errors.New("nil big.Int provided")
	}
	if n.Sign() < 0 {
		return "", errors.New("negative big.Int cannot be encoded")
	}

	const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	if n.Cmp(big.NewInt(0)) == 0 {
		return "0", nil
	}

	var encoded strings.Builder
	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := new(big.Int)

	temp := new(big.Int).Set(n) // Create a copy to preserve the original value

	for temp.Cmp(zero) > 0 {
		temp.DivMod(temp, base, mod)
		encoded.WriteByte(base62Chars[mod.Int64()])
	}

	// Reverse the string
	runes := []rune(encoded.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes), nil
}

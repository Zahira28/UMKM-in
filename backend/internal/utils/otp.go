package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("gagal menghasilkan bilangan acak untuk kode verifikasi: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

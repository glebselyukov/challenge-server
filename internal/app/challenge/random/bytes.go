package random

import (
	"fmt"
	"math/rand"

	"benchmark/internal/pkg/letters"
	"benchmark/internal/pkg/sizes"
)

func RandASCIIBytes(n sizes.ByteSize) ([]byte, error) {
	if n > sizes.MaximumRange {
		return nil, fmt.Errorf("out of range, maximum range: %v\n", sizes.MaximumRange)
	}
	output := make([]byte, n, n)
	randomness := make([]byte, n, n)
	_, err := rand.Read(randomness)
	if err != nil {
		return nil, err
	}
	l := len(letters.LetterBytes)
	for pos := range output {
		random := uint8(randomness[pos])
		randomPos := random % uint8(l)
		output[pos] = letters.LetterBytes[randomPos]
	}
	return output, nil
}

package files

import (
	"os"

	"github.com/prospik/challenge-server/internal/pkg/letters"
	"github.com/prospik/challenge-server/internal/pkg/sizes"
)

func BytesFromData() ([]byte, error) {
	f, err := os.Open(letters.Path)
	if err != nil {
		return nil, err
	}

	bytes := make([]byte, sizes.DefaultFileSize, sizes.DefaultFileSize)
	_, err = f.Read(bytes)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return bytes, nil
}

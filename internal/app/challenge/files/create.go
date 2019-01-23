package files

import (
	"os"
	"strings"

	"github.com/prospik/challenge-server/internal/app/challenge/random"
	"github.com/prospik/challenge-server/internal/pkg/letters"
	"github.com/prospik/challenge-server/internal/pkg/sizes"
)

func CreateDumpData(assets, file string) {
	letters.Path = strings.Join([]string{assets, file}, "/")
	path := letters.Path
	if !fileExists(path) {
		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		randomBytes, err := random.RandASCIIBytes(sizes.DefaultFileSize)
		_, err = f.Write(randomBytes)
		if err != nil {
			panic("can't create file")
		}
	}
}

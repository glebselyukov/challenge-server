package main

import (
	"flag"
	"os"

	"github.com/prospik/challenge-server/internal/app/challenge/files"
	"github.com/prospik/challenge-server/internal/pkg/letters"
	"github.com/prospik/challenge-server/web/api"
)

var (
	port int
	fileName,
	assetsPath string
)

func main() {
	flag.StringVar(&assetsPath, "path", letters.DefaultAssetsPathName, "")
	flag.StringVar(&fileName, "file", letters.DefaultFileName, "")
	flag.IntVar(&port, "port", letters.DefaultPort, "")
	flag.Parse()

	if env, isExist := os.LookupEnv("ASSETS_PATH_NAME"); isExist {
		assetsPath = env
	}

	if env, isExist := os.LookupEnv("FILE_NAME"); isExist {
		fileName = env
	}

	go files.CreateDumpData(assetsPath, fileName)

	api.New(port)
}

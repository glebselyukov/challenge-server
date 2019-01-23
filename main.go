package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"unsafe"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type ByteSize int

const (
	_           = iota
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
)

const (
	letterBytes           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defaultAssetsPathName = "assets"
	defaultFileName       = "data"
	defaultFileSize       = KB * 64
	maximumRange          = GB
)

var (
	path,
	fileName,
	assetsPath string
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

type HashResponse struct {
	Iterations *int64  `json:"iterations"`
	Result     *string `json:"result"`
	Error      error   `json:"error,omitempty"`
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func randomASCIIBytes(n ByteSize) ([]byte, error) {
	if n > maximumRange {
		return nil, fmt.Errorf("out of range, maximum range: %v\n", maximumRange)
	}
	output := make([]byte, n, n)
	randomness := make([]byte, n, n)
	_, err := rand.Read(randomness)
	if err != nil {
		return nil, err
	}
	l := len(letterBytes)
	for pos := range output {
		random := uint8(randomness[pos])
		randomPos := random % uint8(l)
		output[pos] = letterBytes[randomPos]
	}
	return output, nil
}

func doJSONWrite(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

func values(ctx *fasthttp.RequestCtx) {
	response := &HashResponse{}
	n := ctx.QueryArgs().Peek("n")

	var iterations int64
	if len(n) > 0 {
		iterationsParse, err := strconv.ParseInt(b2s(n), 10, 64)
		if err != nil {
			response.Error = err
			doJSONWrite(ctx, 200, response)
			return
		}
		iterations = iterationsParse
	}

	response.Iterations = &iterations

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	bytes := make([]byte, defaultFileSize, defaultFileSize)
	_, err = f.Read(bytes)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	sha := sha256.Sum256(bytes)
	for i := 1; i < int(iterations); i++ {
		sha = sha256.Sum256(sha[:])
	}

	result := base64.URLEncoding.EncodeToString(sha[:])
	response.Result = &result

	doJSONWrite(ctx, 200, response)
}

func createDumpData() {
	if !fileExists(assetsPath) {
		err := os.MkdirAll(assetsPath, os.ModePerm)
		if err != nil {
			panic("can't create dirs")
		}
	}

	path = strings.Join([]string{assetsPath, fileName}, "/")
	if !fileExists(path) {
		f, err := os.Create(path)
		if err != nil {
			panic("can't create file")
		}

		defer f.Close()

		randomBytes, err := randomASCIIBytes(defaultFileSize)
		_, err = f.Write(randomBytes)
		if err != nil {
			panic("can't create file")
		}
	}
}

func main() {
	flag.StringVar(&assetsPath, "path", defaultAssetsPathName, "")
	flag.StringVar(&fileName, "file", defaultFileName, "")
	flag.Parse()

	if env, isExist := os.LookupEnv("ASSETS_PATH_NAME"); isExist {
		assetsPath = env
	}

	if env, isExist := os.LookupEnv("FILE_NAME"); isExist {
		fileName = env
	}

	go createDumpData()

	router := fasthttprouter.New()
	router.GET("/values", values)

	log.Fatal(fasthttp.ListenAndServe(":50000", router.Handler))
}

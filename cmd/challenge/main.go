package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"unsafe"

	"benchmark/internal/app/challenge/files"
	"benchmark/internal/pkg/letters"
	"benchmark/internal/pkg/sizes"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
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

	bytes := make([]byte, sizes.DefaultFileSize, sizes.DefaultFileSize)
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

func main() {
	flag.StringVar(&assetsPath, "path", letters.DefaultAssetsPathName, "")
	flag.StringVar(&fileName, "file", letters.DefaultFileName, "")
	flag.Parse()

	if env, isExist := os.LookupEnv("ASSETS_PATH_NAME"); isExist {
		assetsPath = env
	}

	if env, isExist := os.LookupEnv("FILE_NAME"); isExist {
		fileName = env
	}

	go files.CreateDumpData(assetsPath, fileName)

	router := fasthttprouter.New()
	router.GET("/values", values)

	log.Fatal(fasthttp.ListenAndServe(":50000", router.Handler))
}

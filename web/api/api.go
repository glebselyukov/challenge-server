package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"unsafe"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/prospik/challenge-server/internal/app/challenge/files"
)

type HashResponse struct {
	// Iterations *int64  `json:"iterations"`
	Result *string `json:"result"`
	Error  error   `json:"error,omitempty"`
}

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

func New(port int) {
	router := fasthttprouter.New()
	router.GET("/values", values)

	p := fmt.Sprintf(":%v", port)
	log.Fatal(fasthttp.ListenAndServe(p, router.Handler))
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

	// response.Iterations = &iterations

	bytes, err := files.BytesFromData()
	if err != nil {
		panic(err)
	}

	sha := sha256.Sum256(bytes)
	for i := 1; i < int(iterations); i++ {
		sha = sha256.Sum256(sha[:])
	}

	result := base64.URLEncoding.EncodeToString(sha[:])
	response.Result = &result

	doJSONWrite(ctx, 200, response)
}

package api

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"unsafe"

	"github.com/buaazp/fasthttprouter"
	"github.com/json-iterator/go"
	"github.com/valyala/fasthttp"

	"github.com/prospik/challenge-server/internal/app/challenge/files"
)

type HashResponse struct {
	Result *string `json:"result"`
	Error  error   `json:"error,omitempty"`
}

var (
	strApplicationJSON = []byte("application/json")
)

func New(port int) {
	router := fasthttprouter.New()
	router.GET("/api/values", values)

	p := fmt.Sprintf(":%v", port)
	log.Fatal(fasthttp.ListenAndServe(p, router.Handler))
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func doJSONWrite(ctx *fasthttp.RequestCtx, code int, response interface{}) {
	ctx.Response.Header.SetContentTypeBytes(strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if bytes, err := json.Marshal(response); err == nil {
		ctx.SetBody(bytes)
		return
	}
	ctx.Error("Error", fasthttp.StatusInternalServerError)

}

func values(ctx *fasthttp.RequestCtx) {
	response := &HashResponse{}
	n := ctx.QueryArgs().Peek("n")

	var iterations int
	if len(n) > 0 {
		iterationsParse, err := strconv.Atoi(b2s(n))
		if err != nil {
			response.Error = err
			doJSONWrite(ctx, 200, response)
			return
		}
		iterations = iterationsParse
	}

	bytes, err := files.BytesFromData()
	if err != nil {
		panic(err)
	}

	sha := sha256.Sum256(bytes)
	for i := 1; i < iterations; i++ {
		sha = sha256.Sum256(sha[:])
	}

	base := base64.URLEncoding
	buf := make([]byte, base.EncodedLen(len(sha[:])))
	base.Encode(buf, sha[:])
	result := b2s(buf)
	response.Result = &result

	doJSONWrite(ctx, 200, response)
}

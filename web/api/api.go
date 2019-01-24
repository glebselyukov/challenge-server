package api

import (
	"encoding/base64"
	"strconv"
	"sync"
	"unsafe"

	"github.com/json-iterator/go"
	"github.com/minio/sha256-simd"
	"github.com/valyala/fasthttp"

	"github.com/prospik/challenge-server/internal/app/challenge/files"
)

type hashResponse struct {
	Result *string `json:"result"`
	Error  error   `json:"error,omitempty"`
}

var (
	bytes []byte
	pool  = &sync.Pool{
		New: func() interface{} {
			return &hashResponse{}
		},
	}
)

func New(port int) {
	p := strconv.FormatInt(int64(port), 10)
	addr := ":" + p
	server := &fasthttp.Server{
		Handler:                            values,
		Name:                               "",
		GetOnly:                            true,
		LogAllErrors:                       false,
		DisableHeaderNamesNormalizing:      true,
		SleepWhenConcurrencyLimitsExceeded: 0,
		NoDefaultServerHeader:              true,
	}

	var err error
	bytes, err = files.BytesFromData()
	if err != nil {
		panic(err)
	}

	err = server.ListenAndServe(addr)
	if err != nil {
		panic(err)
	}
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func doJSONWrite(ctx *fasthttp.RequestCtx, code int, response interface{}) {
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetStatusCode(code)
	var json = jsoniter.ConfigFastest
	if bytes, err := json.Marshal(response); err == nil {
		ctx.Response.SetBodyString(b2s(bytes))
		return
	}
	ctx.Error("Error", fasthttp.StatusInternalServerError)

}

func values(ctx *fasthttp.RequestCtx) {
	response := pool.Get().(*hashResponse)

	defer func() {
		response.Result = nil
		response.Error = nil
		pool.Put(response)
	}()

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

	sha := sha256.Sum256(bytes)
	for i := 1; i < iterations; i++ {
		sha = sha256.Sum256(sha[:])
	}

	base := base64.URLEncoding
	buf := make([]byte, base.EncodedLen(sha256.Size))
	base.Encode(buf, sha[:])
	result := b2s(buf)
	response.Result = &result

	doJSONWrite(ctx, 200, response)
}

package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type AccessLog struct {
	Method   string
	Url      string
	Duration string
	ReqBody  string
	RespBody string
	Status   int
}

type ResponseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

type MiddlewareBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewMiddlewareBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc: fn,
	}
}

func (mb *MiddlewareBuilder) AllowReqBody() *MiddlewareBuilder {
	mb.allowReqBody = true
	return mb
}

func (mb *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	mb.allowRespBody = true
	return mb
}

func (mb *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}
		if mb.allowReqBody && ctx.Request.Body != nil {
			body, _ := ctx.GetRawData()
			reader := io.NopCloser(bytes.NewReader(body))
			ctx.Request.Body = reader

			if len(body) > 1024 {
				body = body[:1024]
			}
			al.ReqBody = string(body)
		}

		if mb.allowRespBody {
			ctx.Writer = &ResponseWriter{
				al:             al,
				ResponseWriter: ctx.Writer,
			}
		}

		defer func() {
			al.Duration = time.Since(start).String()
			mb.loggerFunc(ctx, al)
		}()

		ctx.Next()
	}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.al.Status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.al.RespBody = string(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *ResponseWriter) WriteString(data string) (int, error) {
	rw.al.RespBody = data
	return rw.ResponseWriter.WriteString(data)
}

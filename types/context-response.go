package types

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MetaDiv-AI/metadev/internal/mime"
	"github.com/MetaDiv-AI/metaorm"
	"github.com/gin-gonic/gin"
)

func NewResponseContext[ResponseType any]() *responseContext[ResponseType] {
	return &responseContext[ResponseType]{
		startAt: time.Now(),
	}
}

type ResponseContext[ResponseType any] interface {
	OK(data *ResponseType, page ...*metaorm.Pagination)
	Error(err error)
	File(file []byte, filename string)
	Download(file []byte, filename string)
	String(data string)

	MakeGinResponse(ctx *gin.Context) (ok bool)
	SetNewToken(token string)
}

type responseContext[ResponseType any] struct {
	startAt time.Time

	responseType string // json, file, download
	httpStatus   int
	response     *Response[ResponseType]
	raw          string
	filename     string
	file         []byte

	newToken string
}

// OK creates a success response
func (b *responseContext[ResponseType]) OK(data *ResponseType, page ...*metaorm.Pagination) {
	resp := &Response[ResponseType]{
		Success:  true,
		Time:     b.startAt.Format("2006-01-02 15:04:05"),
		Duration: time.Since(b.startAt).Milliseconds(),
		Data:     data,
	}
	if len(page) > 0 {
		resp.Page = page[0]
	}
	raw, err := json.Marshal(resp)
	if err != nil {
		return
	}
	b.responseType = "json"
	b.httpStatus = http.StatusOK
	b.response = resp
	b.raw = string(raw)
}

// Error creates an error response
func (b *responseContext[ResponseType]) Error(err error) {
	resp := &Response[ResponseType]{
		Success:  false,
		Time:     b.startAt.Format("2006-01-02 15:04:05"),
		Duration: time.Since(b.startAt).Milliseconds(),
		Error:    err.Error(),
	}
	raw, err := json.Marshal(resp)
	if err != nil {
		return
	}
	b.responseType = "json"
	b.httpStatus = http.StatusBadRequest
	b.response = resp
	b.raw = string(raw)
}

func (b *responseContext[ResponseType]) File(file []byte, filename string) {
	b.responseType = "file"
	b.httpStatus = http.StatusOK
	b.file = file
	b.filename = filename
}

func (b *responseContext[ResponseType]) Download(file []byte, filename string) {
	b.responseType = "download"
	b.httpStatus = http.StatusOK
	b.file = file
	b.filename = filename
}

// String creates a string response
func (b *responseContext[ResponseType]) String(data string) {
	b.responseType = "string"
	b.httpStatus = http.StatusOK
	b.raw = data
}

func (b *responseContext[ResponseType]) SetNewToken(token string) {
	b.newToken = token
}

func (b *responseContext[ResponseType]) MakeGinResponse(ctx *gin.Context) bool {
	switch b.responseType {
	case "json":
		ctx.JSON(b.httpStatus, b.response)
		return true
	case "string":
		ctx.String(b.httpStatus, b.raw)
		return true
	case "file":
		ctx.Header("Content-Disposition", "filename="+b.filename)
		mimeConvertor := mime.NewConvertor()
		ctx.Data(http.StatusOK, mimeConvertor.FileNameToMime(b.filename), b.file)
		return true
	case "download":
		ctx.Header("Content-Disposition", "attachment; filename="+b.filename)
		ctx.Data(http.StatusOK, "application/octet-stream", b.file)
		return true
	default:
		return false
	}
}

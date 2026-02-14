package types

import (
	"bytes"
	"fmt"
	"io"

	"github.com/MetaDiv-AI/metadev/internal/ginreq"

	"github.com/gin-gonic/gin"
)

// NewRequestContext creates a new request context
func NewRequestContext[RequestType any](ctx *gin.Context) (*requestContext[RequestType], error) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return &requestContext[RequestType]{}, fmt.Errorf("failed to read request body: %s", err.Error())
	}

	// Restore the body for subsequent reads
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Parse request with error handling
	ginExtractor := ginreq.NewExtractor[RequestType](ctx)
	request, err := ginExtractor.ExtractRequest()
	if err != nil {
		return &requestContext[RequestType]{}, err
	}

	// Extract file with error handling
	fileContent, filename, fileSize, err := ginExtractor.ExtractFile()
	if err != nil && ctx.ContentType() == "multipart/form-data" {
		fmt.Println("failed to extract file from multipart request", err.Error())
	}

	return &requestContext[RequestType]{
		request:  request,
		raw:      string(body),
		file:     fileContent,
		filename: filename,
		fileSize: fileSize,
	}, nil
}

type RequestContext[RequestType any] interface {
	// Request returns the request
	Request() *RequestType
	// RawRequest returns the raw request
	RawRequest() string
	// FileRequest returns the file request
	FileRequest() (file []byte, filename string, size int64)
}

type requestContext[RequestType any] struct {
	request  *RequestType
	raw      string
	filename string
	fileSize int64
	file     []byte
}

func (c *requestContext[RequestType]) Request() *RequestType {
	return c.request
}

func (c *requestContext[RequestType]) RawRequest() string {
	return c.raw
}

func (c *requestContext[RequestType]) FileRequest() (file []byte, filename string, size int64) {
	return c.file, c.filename, c.fileSize
}

package types

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func NewGinContext(ctx *gin.Context) *ginContext {
	return &ginContext{ctx: ctx}
}

type GinContext interface {
	// Gin returns the gin context
	Gin() *gin.Context
	// IP returns the client IP
	IP() string
	// Agent returns the user agent
	Agent() string
	// Method returns the request method
	Method() string
	// Path returns the request path
	Path() string
	// Locale returns the locale
	Locale() string
	// Header returns the header
	Header(header string) string
	// Bearer returns the bearer token
	Bearer() string
}

type ginContext struct {
	ctx *gin.Context
}

func (c *ginContext) Gin() *gin.Context {
	return c.ctx
}

func (c *ginContext) IP() string {
	return c.ctx.ClientIP()
}

func (c *ginContext) Agent() string {
	return c.ctx.GetHeader("User-Agent")
}

func (c *ginContext) Method() string {
	return c.ctx.Request.Method
}

func (c *ginContext) Path() string {
	return c.ctx.FullPath()
}

func (g *ginContext) Locale() string {
	locale := g.Header("X-Locale")
	if locale == "" {
		return "en"
	}
	return locale
}

func (c *ginContext) Header(header string) string {
	return c.ctx.GetHeader(header)
}

func (c *ginContext) Bearer() string {
	token := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
		c.ctx.GetHeader("Authorization"), "Bearer", ""), "BEARER", "bearer"), "bearer", ""), " ", "")
	if token == "" {
		token = c.ctx.Query("token")
	}
	return token
}

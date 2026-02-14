package types

import (
	"github.com/gin-gonic/gin"
)

func NewJwtContext(gin *gin.Context) *jwtContext {
	return &jwtContext{ginContext: ginContext{ctx: gin}}
}

type JwtContext interface {
	GinContext
	// Jwt returns the JWT
	Jwt() Jwt
	// IsAdmin returns if the user is an admin
	IsAdmin() bool
	// UserId returns the user ID
	UserId() uint
	// WorkspaceId returns the workspace ID
	WorkspaceId() uint
	// WorkspaceMemberId returns the workspace member ID
	WorkspaceMemberId() uint
}

type jwtContext struct {
	ginContext
	jwt Jwt
}

func (c *jwtContext) Jwt() Jwt {
	if c.jwt != nil {
		return c.jwt
	}
	token := c.ginContext.Bearer()
	if token == "" {
		return nil
	}
	jwt, err := ParseJwt(token)
	if err != nil {
		return nil
	}
	c.jwt = jwt
	return c.jwt
}

func (c *jwtContext) IsAdmin() bool {
	jwt := c.Jwt()
	if jwt == nil {
		return false
	}
	return jwt.IsAdmin()
}

func (c *jwtContext) UserId() uint {
	jwt := c.Jwt()
	if jwt == nil {
		return 0
	}
	return jwt.UserId()
}

func (c *jwtContext) WorkspaceId() uint {
	jwt := c.Jwt()
	if jwt == nil {
		return 0
	}
	return jwt.WorkspaceId()
}

func (c *jwtContext) WorkspaceMemberId() uint {
	jwt := c.Jwt()
	if jwt == nil {
		return 0
	}
	return jwt.WorkspaceMemberId()
}

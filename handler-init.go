package metadev

import (
	"github.com/MetaDiv-AI/metamongo"

	"github.com/MetaDiv-AI/metaorm"

	"github.com/MetaDiv-AI/metadev/types"
)

// InitFunc creates a new init handler builder
func InitFunc(module types.Module) *initNameBuilder {
	return &initNameBuilder{module: module}
}

type initNameBuilder struct {
	module types.Module
}

// Name sets the name of the init handler
func (b *initNameBuilder) Name(name string) *initHandlerBuilder {
	return &initHandlerBuilder{module: b.module, name: name}
}

type initHandlerBuilder struct {
	module types.Module
	name   string
}

// Handler sets the handler of the init handler
func (b *initHandlerBuilder) Handler(handler func(db metaorm.Database, mongo metamongo.Database, logger types.Logger)) types.InitHandler {
	h := types.NewInitHandler(b.module, b.name, handler)
	b.module.RegisterInitHandler(h)
	return h
}

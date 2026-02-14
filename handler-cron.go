package metadev

import (
	"github.com/MetaDiv-AI/metadev/types"
	"github.com/MetaDiv-AI/metamongo"
	"github.com/MetaDiv-AI/metaorm"
)

func NewCronHandler(module types.Module) *cronNameBuilder {
	return &cronNameBuilder{module: module}
}

type cronNameBuilder struct {
	module types.Module
}

func (b *cronNameBuilder) Name(name string) *cronSpecBuilder {
	return &cronSpecBuilder{module: b.module, name: name}
}

type cronSpecBuilder struct {
	module types.Module
	name   string
}

func (b *cronSpecBuilder) Spec(spec string) *cronHandlerBuilder {
	return &cronHandlerBuilder{module: b.module, name: b.name, spec: spec}
}

type cronHandlerBuilder struct {
	module types.Module
	name   string
	spec   string
}

func (b *cronHandlerBuilder) Handler(handler func(db metaorm.Database, mongo metamongo.Database, logger types.Logger)) types.CronHandler {
	h := types.NewCronHandler(b.module, b.name, b.spec, handler)
	b.module.RegisterCronHandler(h)
	return h
}

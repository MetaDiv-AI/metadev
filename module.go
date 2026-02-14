package metadev

import "github.com/MetaDiv-AI/metadev/types"

// NewModule creates a new module builder
func NewModule(app types.App) *moduleNameBuilder {
	return &moduleNameBuilder{app: app}
}

type moduleNameBuilder struct {
	app types.App
}

// Name sets the name of the module
func (b *moduleNameBuilder) Name(name string) types.Module {
	m := types.NewModule(b.app, name)
	b.app.RegisterModule(m)
	return m
}

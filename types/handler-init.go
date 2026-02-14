package types

import (
	"github.com/MetaDiv-AI/metamongo"

	"github.com/MetaDiv-AI/metaorm"
)

func NewInitHandler(module Module, name string, handler func(db metaorm.Database, mongo metamongo.Database, logger Logger)) InitHandler {
	return &initHandler{module: module, name: name, handler: handler}
}

type InitHandler interface {
	// Module returns the module of the init handler
	Module() Module
	// Name returns the name of the init handler
	Name() string
	// Handler returns the handler of the init handler
	Handler() func(db metaorm.Database, mongo metamongo.Database, logger Logger)
}

type initHandler struct {
	module  Module
	name    string
	handler func(db metaorm.Database, mongo metamongo.Database, logger Logger)
}

func (h *initHandler) Module() Module {
	return h.module
}

func (h *initHandler) Name() string {
	return h.name
}

func (h *initHandler) Handler() func(db metaorm.Database, mongo metamongo.Database, logger Logger) {
	return h.handler
}

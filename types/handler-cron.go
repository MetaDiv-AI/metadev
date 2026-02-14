package types

import (
	"github.com/MetaDiv-AI/metamongo"

	"github.com/MetaDiv-AI/metaorm"
)

func NewCronHandler(module Module, name string, spec string, handler func(db metaorm.Database, mongo metamongo.Database, logger Logger)) CronHandler {
	return &cronHandler{module: module, name: name, spec: spec, handler: handler}
}

type CronHandler interface {
	// Module returns the module of the cron handler
	Module() Module
	// Name returns the name of the cron handler
	Name() string
	// Spec returns the spec of the cron handler
	Spec() string
	// Handler returns the handler of the cron handler
	Handler() func(db metaorm.Database, mongo metamongo.Database, logger Logger)
}

type cronHandler struct {
	module  Module
	name    string
	spec    string
	handler func(db metaorm.Database, mongo metamongo.Database, logger Logger)
}

func (h *cronHandler) Module() Module {
	return h.module
}

func (h *cronHandler) Name() string {
	return h.name
}

func (h *cronHandler) Spec() string {
	return h.spec
}

func (h *cronHandler) Handler() func(db metaorm.Database, mongo metamongo.Database, logger Logger) {
	return h.handler
}

package types

import (
	"fmt"
	"sort"
)

func NewModule(app App, name string) Module {
	return &ModuleImpl{
		app:          app,
		name:         name,
		apiHandlers:  make(map[string]ApiHandler),
		cronHandlers: make(map[string]CronHandler),
		initHandlers: make(map[string]InitHandler),
		wsHandlers:   make(map[string]WsHandler),
	}
}

type Module interface {
	App() App
	Name() string
	ApiHandlers() []ApiHandler
	CronHandlers() []CronHandler
	InitHandlers() []InitHandler
	WsHandlers() []WsHandler
	RegisterApiHandler(handler ApiHandler)
	RegisterCronHandler(handler CronHandler)
	RegisterInitHandler(handler InitHandler)
	RegisterWsHandler(handler WsHandler)
}

type ModuleImpl struct {
	app          App
	name         string
	apiHandlers  map[string]ApiHandler
	cronHandlers map[string]CronHandler
	initHandlers map[string]InitHandler
	wsHandlers   map[string]WsHandler
}

func (m *ModuleImpl) App() App {
	return m.app
}

func (m *ModuleImpl) Name() string {
	return m.name
}

func (m *ModuleImpl) ApiHandlers() []ApiHandler {
	keys := make([]string, 0, len(m.apiHandlers))
	for key := range m.apiHandlers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	handlers := make([]ApiHandler, 0, len(m.apiHandlers))
	for _, key := range keys {
		handlers = append(handlers, m.apiHandlers[key])
	}
	return handlers
}

func (m *ModuleImpl) CronHandlers() []CronHandler {
	keys := make([]string, 0, len(m.cronHandlers))
	for key := range m.cronHandlers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	handlers := make([]CronHandler, 0, len(m.cronHandlers))
	for _, key := range keys {
		handlers = append(handlers, m.cronHandlers[key])
	}
	return handlers
}

func (m *ModuleImpl) InitHandlers() []InitHandler {
	keys := make([]string, 0, len(m.initHandlers))
	for key := range m.initHandlers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	handlers := make([]InitHandler, 0, len(m.initHandlers))
	for _, key := range keys {
		handlers = append(handlers, m.initHandlers[key])
	}
	return handlers
}

func (m *ModuleImpl) WsHandlers() []WsHandler {
	keys := make([]string, 0, len(m.wsHandlers))
	for key := range m.wsHandlers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	handlers := make([]WsHandler, 0, len(m.wsHandlers))
	for _, key := range keys {
		handlers = append(handlers, m.wsHandlers[key])
	}
	return handlers
}

func (m *ModuleImpl) RegisterApiHandler(handler ApiHandler) {
	_, ok := m.apiHandlers[handler.Name()]
	if ok {
		panic(fmt.Sprintf("api handler %s already registered for module %s", handler.Name(), m.Name()))
	}
	m.apiHandlers[handler.Name()] = handler
}

func (m *ModuleImpl) RegisterCronHandler(handler CronHandler) {
	_, ok := m.cronHandlers[handler.Name()]
	if ok {
		panic(fmt.Sprintf("cron handler %s already registered for module %s", handler.Name(), m.Name()))
	}
	m.cronHandlers[handler.Name()] = handler
}

func (m *ModuleImpl) RegisterInitHandler(handler InitHandler) {
	_, ok := m.initHandlers[handler.Name()]
	if ok {
		panic(fmt.Sprintf("init handler %s already registered for module %s", handler.Name(), m.Name()))
	}
	m.initHandlers[handler.Name()] = handler
}

func (m *ModuleImpl) RegisterWsHandler(handler WsHandler) {
	_, ok := m.wsHandlers[handler.Name()]
	if ok {
		panic(fmt.Sprintf("ws handler %s already registered for module %s", handler.Name(), m.Name()))
	}
	m.wsHandlers[handler.Name()] = handler
}

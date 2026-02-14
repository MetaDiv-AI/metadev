package types

import (
	"sync"
	"time"

	"github.com/MetaDiv-AI/metaorm"

	"github.com/gorilla/websocket"
)

type Values struct {
	Values map[string]any `json:"values"`
}

func (v *Values) SetString(key string, value string) {
	v.Values[key] = value
}

func (v *Values) SetInt(key string, value int) {
	v.Values[key] = value
}

func (v *Values) SetFloat64(key string, value float64) {
	v.Values[key] = value
}

func (v *Values) SetBool(key string, value bool) {
	v.Values[key] = value
}

func (v *Values) String(key string) string {
	return v.Values[key].(string)
}

func (v *Values) Int(key string) int {
	return v.Values[key].(int)
}

func (v *Values) Float64(key string) float64 {
	return v.Values[key].(float64)
}

func (v *Values) Bool(key string) bool {
	b, ok := v.Values[key].(bool)
	if !ok {
		return false
	}
	return b
}

func NewWsConnContext[ResponseType any](wsConn *websocket.Conn, writeMutex *sync.Mutex) *wsConnContext[ResponseType] {
	return &wsConnContext[ResponseType]{
		wsConn:     wsConn,
		values:     &Values{Values: make(map[string]any)},
		writeMutex: writeMutex,
	}
}

type WsConnContext[ResponseType any] interface {
	WsConn() *websocket.Conn
	SendMessage(action string, data *ResponseType, page ...*metaorm.Pagination) error
	SendError(err error) error
	Values() *Values
}

type wsConnContext[ResponseType any] struct {
	wsConn     *websocket.Conn
	values     *Values
	writeMutex *sync.Mutex
}

func (c *wsConnContext[ResponseType]) WsConn() *websocket.Conn {
	return c.wsConn
}

func (c *wsConnContext[ResponseType]) Values() *Values {
	return c.values
}

func (c *wsConnContext[ResponseType]) SendMessage(action string, data *ResponseType, page ...*metaorm.Pagination) error {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()

	var pagination *metaorm.Pagination
	if len(page) > 0 {
		pagination = page[0]
	}
	return c.WsConn().WriteJSON(WsResponse[ResponseType]{
		Action: action,
		Response: Response[ResponseType]{
			Success: true,
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			Data:    data,
			Page:    pagination,
		},
	})
}

func (c *wsConnContext[ResponseType]) SendError(err error) error {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()

	return c.WsConn().WriteJSON(WsResponse[ResponseType]{
		Action: "error",
		Response: Response[ResponseType]{
			Success: false,
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			Error:   err.Error(),
		},
	})
}

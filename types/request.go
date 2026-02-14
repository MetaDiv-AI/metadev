package types

type WsMessage[RequestType any] struct {
	Action  string       `json:"action"`
	Request *RequestType `json:"request"`
}

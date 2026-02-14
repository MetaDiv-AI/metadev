package types

import "github.com/MetaDiv-AI/metaorm"

type Response[ResponseType any] struct {
	Success bool `json:"success"`

	Time     string `json:"time"`
	Duration int64  `json:"duration"` // in milliseconds

	Data  *ResponseType       `json:"data,omitempty"`
	Error string              `json:"error,omitempty"`
	Page  *metaorm.Pagination `json:"page,omitempty"`
}

type WsResponse[ResponseType any] struct {
	Action string `json:"action"`
	Response[ResponseType]
}

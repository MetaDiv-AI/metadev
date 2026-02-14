package typescript

import "reflect"

type ApiInfo struct {
	Name     string
	Route    string
	Method   string
	Uris     []string
	Forms    []string
	Request  string
	Response string

	RequestType  reflect.Type
	ResponseType reflect.Type
}

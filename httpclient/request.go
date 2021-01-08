package httpclient

import "github.com/zfd81/rooster/types/container"

type Request struct {
	Header Header
	Params container.JsonMap
}

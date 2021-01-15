package script

import (
	"github.com/zfd81/rock/core"
	"github.com/zfd81/rock/script/otto"
)

func New() core.Script {
	return otto.New()
}

func NewWithContext(ctx core.Context) core.Script {
	return otto.NewWithContext(ctx)
}

func NewWithProcessor(processor core.Processor) core.Script {
	return otto.NewWithProcessor(processor)
}

func GetSdk() string {
	return otto.GetSdk()
}

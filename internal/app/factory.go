package app

import (
	"github.com/VadimGossip/gitPatchTool/internal/walker"
)

type Factory struct {
	walkerService service.Service
}

var factory *Factory

func newFactory() (*Factory, error) {
	factory = &Factory{}
	factory.walkerService = service.NewService()
	return factory, nil
}

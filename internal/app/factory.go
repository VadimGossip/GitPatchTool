package app

import (
	"github.com/VadimGossip/gitPatchTool/internal/filewalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool"
)

type Factory struct {
	walkerService  filewalker.Service
	oraToolService oratool.Service
}

var factory *Factory

func newFactory() (*Factory, error) {
	factory = &Factory{}
	factory.walkerService = filewalker.NewService()
	factory.oraToolService = oratool.NewService(factory.walkerService)
	return factory, nil
}

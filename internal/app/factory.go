package app

import (
	"github.com/VadimGossip/gitPatchTool/internal/filewalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/patcher"
)

type Factory struct {
	walkerService filewalker.Service

	oraToolExtractor extractor.Service
	oraToolService   patcher.Service
}

var factory *Factory

func newFactory() (*Factory, error) {
	factory = &Factory{}
	factory.walkerService = filewalker.NewService()
	factory.oraToolExtractor = extractor.NewService()
	factory.oraToolService = patcher.NewService(factory.walkerService, factory.oraToolExtractor)
	return factory, nil
}

package app

import (
	"github.com/VadimGossip/gitPatchTool/internal/filewalker"
	"github.com/VadimGossip/gitPatchTool/internal/gitwalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/patcher"
)

type Factory struct {
	dbAdapter *DBAdapter

	gitWalkerService  gitwalker.Service
	fileWalkerService filewalker.Service

	oraToolExtractor extractor.Service
	oraToolService   patcher.Service
}

var factory *Factory

func newFactory(dbAdapter *DBAdapter) (*Factory, error) {
	factory = &Factory{dbAdapter: dbAdapter}
	factory.gitWalkerService = gitwalker.NewService(dbAdapter.gitWalkerRepo)
	factory.fileWalkerService = filewalker.NewService()
	factory.oraToolExtractor = extractor.NewService()
	factory.oraToolService = patcher.NewService(factory.gitWalkerService, factory.fileWalkerService, factory.oraToolExtractor)
	return factory, nil
}

package app

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/filewalker"
	"github.com/VadimGossip/gitPatchTool/internal/gitwalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/patcher"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/splitter"
)

type Factory struct {
	dbAdapter *DBAdapter

	gitWalkerService  gitwalker.Service
	fileWalkerService filewalker.Service

	oraToolExtractor extractor.Service
	oraToolPatcher   patcher.Service
	oraToolSplitter  splitter.Service
}

var factory *Factory

func newFactory(cfg *domain.Config, dbAdapter *DBAdapter) (*Factory, error) {
	factory = &Factory{dbAdapter: dbAdapter}
	factory.gitWalkerService = gitwalker.NewService(dbAdapter.gitWalkerRepo)
	factory.fileWalkerService = filewalker.NewService()
	factory.oraToolExtractor = extractor.NewService(factory.fileWalkerService)
	factory.oraToolSplitter = splitter.NewService(cfg, factory.oraToolExtractor)
	factory.oraToolPatcher = patcher.NewService(cfg, factory.gitWalkerService, factory.oraToolExtractor)
	return factory, nil
}

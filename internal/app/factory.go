package app

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/filewalker"
	"github.com/VadimGossip/gitPatchTool/internal/gitwalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/patcher"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/splitter"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/writer"
)

type Factory struct {
	dbAdapter *DBAdapter

	gitWalkerService  gitwalker.Service
	fileWalkerService filewalker.Service

	oraToolExtractor extractor.Service
	oraToolPatcher   patcher.Service
	oraToolSplitter  splitter.Service
	oraToolWriter    writer.Service
}

var factory *Factory

func newFactory(cfg *domain.Config, dbAdapter *DBAdapter) (*Factory, error) {
	factory = &Factory{dbAdapter: dbAdapter}
	factory.gitWalkerService = gitwalker.NewService(dbAdapter.gitWalkerRepo)
	factory.fileWalkerService = filewalker.NewService()
	factory.oraToolExtractor = extractor.NewService(factory.fileWalkerService)
	factory.oraToolWriter = writer.NewService()
	factory.oraToolSplitter = splitter.NewService(cfg, factory.fileWalkerService, factory.oraToolExtractor)
	factory.oraToolPatcher = patcher.NewService(cfg, factory.gitWalkerService, factory.oraToolExtractor, factory.oraToolWriter)
	return factory, nil
}

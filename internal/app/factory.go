package app

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/file"
	"github.com/VadimGossip/gitPatchTool/internal/gitwalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/patcher"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/splitter"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/writer"
)

type Factory struct {
	dbAdapter *DBAdapter

	gitWalkerService gitwalker.Service
	fileService      file.Service

	oraToolExtractor extractor.Service
	oraToolPatcher   patcher.Service
	oraToolSplitter  splitter.Service
	oraToolWriter    writer.Service
}

var factory *Factory

func newFactory(cfg *domain.Config, dbAdapter *DBAdapter) (*Factory, error) {
	factory = &Factory{dbAdapter: dbAdapter}
	factory.gitWalkerService = gitwalker.NewService(dbAdapter.gitWalkerRepo)
	factory.fileService = file.NewService()
	factory.oraToolExtractor = extractor.NewService(factory.fileService)
	factory.oraToolWriter = writer.NewService()
	factory.oraToolSplitter = splitter.NewService(cfg, factory.fileService, factory.oraToolExtractor)
	factory.oraToolPatcher = patcher.NewService(cfg, factory.fileService, factory.gitWalkerService, factory.oraToolExtractor, factory.oraToolWriter)
	return factory, nil
}

package patcher

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/gitwalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/writer"
)

type Service interface {
	CreatePatch() error
}

type service struct {
	cfg       *domain.Config
	gitWalker gitwalker.Service
	extractor extractor.Service
	writer    writer.Service
}

var _ Service = (*service)(nil)

func NewService(cfg *domain.Config, gitWalker gitwalker.Service, extractor extractor.Service, writer writer.Service) *service {
	return &service{cfg: cfg, gitWalker: gitWalker, extractor: extractor, writer: writer}
}

func (s *service) CreatePatch() error {
	gitFiles, err := s.gitWalker.GetFilesChanged(s.cfg.CommitId)
	if err != nil {
		return err
	}

	commitMsg, err := s.gitWalker.FormCurCommitHeaderMsg(s.cfg.CommitId)
	if err != nil {
		return err
	}

	oracleObj := s.extractor.CreateOracleObjects(s.cfg.Path.RootDir, s.cfg.Path.InstallDir, gitFiles)

	return s.writer.CreateInstallFiles(s.cfg.Path.RootDir, s.cfg.Path.InstallDir, commitMsg, oracleObj)
}

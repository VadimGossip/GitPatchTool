package patcher

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/file"
	"github.com/VadimGossip/gitPatchTool/internal/gitwalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/writer"
)

type Service interface {
	CreatePatch() error
}

type service struct {
	cfg       *domain.Config
	file      file.Service
	gitWalker gitwalker.Service
	extractor extractor.Service
	writer    writer.Service
}

var _ Service = (*service)(nil)

func NewService(cfg *domain.Config, file file.Service, gitWalker gitwalker.Service, extractor extractor.Service, writer writer.Service) *service {
	return &service{cfg: cfg, file: file, gitWalker: gitWalker, extractor: extractor, writer: writer}
}

func (s *service) removeSessionFiles() error {
	for _, fName := range []string{domain.ErrorLogFileName, domain.WarningLogFileName} {
		if err := s.file.DeleteFile(s.cfg.Path.InstallDir + fName); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) CreatePatch() error {
	gitFiles, err := s.gitWalker.GetFilesChanged(s.cfg.CommitId)
	if err != nil {
		return err
	}
	oracleFiles := s.extractor.ExtractOracleObjects(gitFiles)

	installFiles := s.writer.CreateInstallLines(s.cfg.Path.InstallDir, oracleFiles)

	if err := s.removeSessionFiles(); err != nil {
		return err
	}

	for _, iFile := range installFiles {
		if err := s.file.CreateFile(iFile.Path, iFile.FileLines); err != nil {
			return err
		}
	}

	return nil
}

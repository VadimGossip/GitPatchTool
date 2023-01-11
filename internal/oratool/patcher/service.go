package patcher

import (
	"fmt"
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

/*есть ощущение, что про файлы патчер не должен знать ничего, он просто получает список оракловых объектов
и передает его создателю скриптов
*/

func (s *service) CreatePatch() error {
	gitFiles, err := s.gitWalker.GetFilesChanged(s.cfg.CommitId)
	if err != nil {
		return err
	}
	oracleFiles := s.extractor.ExtractOracleObjects(gitFiles)

	installFiles := s.writer.CreateInstallLines(oracleFiles)

	for _, iFile := range installFiles {
		for _, line := range iFile.FileLines {
			fmt.Println(line)
		}
	}

	return nil
}

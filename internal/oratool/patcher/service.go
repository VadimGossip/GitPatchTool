package patcher

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/gitwalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/sirupsen/logrus"
)

type Service interface {
	CreatePatch() error
	FixObjectFile() error
}

type service struct {
	cfg       *domain.Config
	gitWalker gitwalker.Service
	extractor extractor.Service
}

var _ Service = (*service)(nil)

func NewService(cfg *domain.Config, gitWalker gitwalker.Service, extractor extractor.Service) *service {
	return &service{cfg: cfg, gitWalker: gitWalker, extractor: extractor}
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

	for _, val := range oracleFiles {
		logrus.Infof("file %+v", val)
	}

	return nil
}

func (s *service) FixObjectFile() error {
	return nil
}

package oratool

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/filewalker"
	"github.com/sirupsen/logrus"
)

type Service interface {
	CreatePatch() error
	FixObjectFile() error
}

type service struct {
	fileWalker filewalker.Service
}

var _ Service = (*service)(nil)

func NewService(fileWalker filewalker.Service) *service {
	return &service{fileWalker: fileWalker}
}

func (s *service) CreatePatch() error {

	files, err := s.fileWalker.Walk("e:\\WorkSpace\\TCS_Oracle\\", domain.OracleFileType)
	if err != nil {
		logrus.Fatalf("Fail to collect files %s", err)
	}
	logrus.Infof("files count %d", len(files))
	return nil
}

func (s *service) FixObjectFile() error {
	return nil
}

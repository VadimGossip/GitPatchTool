package patcher

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/filewalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/sirupsen/logrus"
)

type Service interface {
	CreatePatch() error
	FixObjectFile() error
}

type service struct {
	fileWalker filewalker.Service
	extractor  extractor.Service
}

var _ Service = (*service)(nil)

func NewService(fileWalker filewalker.Service, extractor extractor.Service) *service {
	return &service{fileWalker: fileWalker, extractor: extractor}
}

func (s *service) CreatePatch() error {

	files, err := s.fileWalker.Walk("e:\\WorkSpace\\TCS_Oracle\\", domain.OracleFileType)
	if err != nil {
		logrus.Fatalf("Fail to collect files %s", err)
	}

	oraObjects := s.extractor.ExtractOracleObjects(files)

	for _, obj := range oraObjects {
		if len(obj.Errors) > 0 {
			logrus.Info(obj)
		}
	}

	return nil
}

func (s *service) FixObjectFile() error {
	return nil
}

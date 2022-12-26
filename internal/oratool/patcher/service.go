package patcher

import (
	"github.com/VadimGossip/gitPatchTool/internal/gitwalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/sirupsen/logrus"
)

type Service interface {
	CreatePatch() error
	FixObjectFile() error
}

type service struct {
	gitWalker gitwalker.Service
	extractor extractor.Service
}

var _ Service = (*service)(nil)

func NewService(gitWalker gitwalker.Service, extractor extractor.Service) *service {
	return &service{gitWalker: gitWalker, extractor: extractor}
}

func (s *service) CreatePatch() error {
	gitFiles, err := s.gitWalker.GetFilesChanged("9b7c2a074bfbf22256b3728629182fe9686c9773")
	if err != nil {
		return err
	}
	logrus.Info(gitFiles)
	oracleFiles := s.extractor.ExtractOracleObjects("e:\\WorkSpace\\TCS_Oracle\\", gitFiles)

	logrus.Infof("oracleFiles %+v", oracleFiles)

	//folderFiles, err := s.fileWalker.Walk("e:\\WorkSpace\\TCS_Oracle\\", domain.OracleFileType)
	//if err != nil {
	//	logrus.Fatalf("Fail to collect files %s", err)
	//}

	return nil
}

func (s *service) FixObjectFile() error {
	return nil
}

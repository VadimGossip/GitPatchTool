package splitter

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/sirupsen/logrus"
)

type Service interface {
	SplitTableFiles() error
}

type service struct {
	cfg       *domain.Config
	extractor extractor.Service
}

var _ Service = (*service)(nil)

func NewService(cfg *domain.Config, extractor extractor.Service) *service {
	return &service{cfg: cfg, extractor: extractor}
}

func (s *service) SplitTableFiles() error {
	oraObjects, err := s.extractor.WalkAndExtractOracleObjects(s.cfg.Path.RootDir)
	if err != nil {
		return err
	}

	filteredObj := make([]domain.OracleObject, 0)

	logrus.Info(oraObjects)

	for _, val := range oraObjects {
		if val.ObjectType == domain.OracleTableType {
			filteredObj = append(filteredObj, val)
		}
	}

	logrus.Infof("len filteredObj %d", len(filteredObj))

	return nil
}

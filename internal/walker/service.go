package service

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"os"
	"path/filepath"
)

type Service interface {
	Walk(path string, fileType int) ([]domain.File, error)
}

type service struct {
}

var _ Service = (*service)(nil)

func NewService() *service {
	return &service{}
}

func (s *service) Walk(path string, fileType int) ([]domain.File, error) {
	result := make([]domain.File, 0)
	extMap := make(map[string]struct{})
	if fileType == domain.OracleFileType {
		extMap[".sql"] = struct{}{}
	}
	if err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if _, ok := extMap[filepath.Ext(info.Name())]; ok {
					result = append(result, domain.File{
						Name: info.Name(),
						Path: path,
						Type: fileType,
					})
				}
			}
			return nil
		}); err != nil {
		return nil, err
	}
	return result, nil
}

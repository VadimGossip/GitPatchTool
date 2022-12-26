package filewalker

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"os"
	"path/filepath"
)

type Service interface {
	CheckFileExists(path string) bool
	LostFiles(allList, searchList []domain.File) []domain.File
	Walk(path string, fileType int) ([]domain.File, error)
}

type service struct {
}

var _ Service = (*service)(nil)

func NewService() *service {
	return &service{}
}

func (s *service) LostFiles(allList, searchList []domain.File) []domain.File {
	allMap := make(map[string]struct{})
	result := make([]domain.File, 0)

	for _, file := range allList {
		allMap[file.Path] = struct{}{}
	}

	for _, sFile := range searchList {
		if _, ok := allMap[sFile.Path]; !ok {
			result = append(result, sFile)
		}
	}

	return result
}

func (s *service) CheckFileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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

package file

import (
	"bufio"
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"os"
	"path/filepath"
	"strings"
)

type Service interface {
	CheckFileExists(path string) bool
	SearchStrInFile(starts, path string) (string, error)
	LostFiles(allList, searchList []domain.File) []domain.File
	Walk(path string, extFilter []string) ([]domain.File, error)
	CreateFile(path string, lines []string) error
	DeleteFile(path string) error
	ReadFile(path string) error
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

func (s *service) SearchStrInFile(searchStr, filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(strings.ToLower(scanner.Text()), searchStr) {
			return strings.ToLower(scanner.Text()), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", nil
}

func (s *service) Walk(path string, extFilter []string) ([]domain.File, error) {
	result := make([]domain.File, 0)
	extMap := make(map[string]struct{})
	for _, ext := range extFilter {
		extMap[ext] = struct{}{}
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
					})
				}
			}
			return nil
		}); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) CreateFile(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		_, err = f.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) DeleteFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
	}
	return nil
}

func (s *service) ReadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

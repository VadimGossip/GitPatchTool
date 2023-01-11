package gitwalker

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
)

type Service interface {
	GetFilesChanged(curCommitHash string) ([]domain.File, error)
}

type service struct {
	repo Repository
}

var _ Service = (*service)(nil)

func NewService(repository Repository) *service {
	return &service{repo: repository}
}

func (s *service) leaveLastState(files []domain.File) []domain.File {
	fkMap := make(map[string]int)
	resultMap := make(map[domain.File]struct{})
	for _, val := range files {
		if fkVal, ok := fkMap[val.Name]; ok {
			delete(resultMap, val)
			if (fkVal == domain.AddAction && val.GitAction == domain.ModifyAction) || fkVal == domain.RenameAction {
				val.GitAction = fkVal
			}
			resultMap[val] = struct{}{}
		} else {
			fkMap[val.Name] = val.GitAction
			resultMap[val] = struct{}{}
		}
	}

	result := make([]domain.File, 0, len(fkMap))
	for key := range resultMap {
		result = append(result, key)
	}

	return result
}

func (s *service) GetFilesChanged(curCommitHash string) ([]domain.File, error) {
	fromCommit, err := s.repo.GetPreviousCommit(curCommitHash)
	if err != nil {
		return nil, err
	}

	headCommit, err := s.repo.GetHeadCommit()
	if err != nil {
		return nil, err
	}

	files, err := s.repo.GetFilesDiff(headCommit, fromCommit)
	if err != nil {
		return nil, err
	}

	return s.leaveLastState(files), nil
}

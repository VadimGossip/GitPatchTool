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

/*
  A add -> A del  A del
  A del -> A add  A add
  A add -> A ch   A add
  A ch  -> A del  A del
  A ren A1 -> A del \ A1 add





*/

func (s *service) leaveLastState(files []domain.File) []domain.File {
	fkMap := make(map[string]int)
	resultMap := make(map[domain.File]struct{})
	for _, val := range files {
		key := val
		if fkVal, ok := fkMap[val.Name]; ok {
			key.GitAction = fkVal
			delete(resultMap, key)
			if !(fkVal == domain.AddAction && val.GitAction == domain.ModifyAction) {
				fkMap[val.Name] = val.GitAction
				key.GitAction = val.GitAction
			}
			resultMap[key] = struct{}{}
		} else {
			fkMap[val.Name] = val.GitAction
			if val.GitAction == domain.RenameAction {
				fkMap[val.InitialName] = domain.DeleteAction
				key.Name = val.InitialName
				key.ShortPath = val.ShortPath
				key.Path = val.Path
				key.GitAction = domain.DeleteAction
			}
			resultMap[key] = struct{}{}
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

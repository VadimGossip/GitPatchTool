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
	resultMap := make(map[string]domain.File)
	for _, val := range files {
		if rv, ok := resultMap[val.Path]; ok {
			if (rv.GitDetails.Action == domain.AddAction && val.GitDetails.Action == domain.ModifyAction) || rv.GitDetails.Action == domain.RenameAction {
				if rv.GitDetails.Action == domain.RenameAction {
					val.GitDetails.InitialName = rv.GitDetails.InitialName
					val.GitDetails.InitialPath = rv.GitDetails.InitialPath
				}
				val.GitDetails.Action = rv.GitDetails.Action

			}
		}
		resultMap[val.Path] = val
	}

	result := make([]domain.File, 0, len(resultMap))
	for _, rv := range resultMap {
		result = append(result, rv)
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

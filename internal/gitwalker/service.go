package gitwalker

import (
	"fmt"
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"strings"
)

type Service interface {
	GetFilesChanged(curCommitHash string) ([]domain.File, error)
	FormCurCommitHeaderMsg(curCommitHash string) (string, error)
}

type service struct {
	repo Repository
}

var _ Service = (*service)(nil)

func NewService(repository Repository) *service {
	return &service{repo: repository}
}

//добавить параметр New
func (s *service) leaveLastState(files []domain.File) []domain.File {
	resultMap := make(map[string]domain.File)
	for _, val := range files {
		if rv, ok := resultMap[val.Path]; ok {
			val.GitDetails.New = rv.GitDetails.New
			if rv.GitDetails.Action == domain.RenameAction {
				val.GitDetails.InitialName = rv.GitDetails.InitialName
				val.GitDetails.InitialPath = rv.GitDetails.InitialPath
				delete(resultMap, rv.GitDetails.InitialPath)
			}
		}
		resultMap[val.Path] = val
	}

	result := make([]domain.File, 0, len(resultMap))
	for _, rv := range resultMap {
		if !(rv.GitDetails.New && rv.GitDetails.Action == domain.DeleteAction) {
			if rv.GitDetails.New {
				rv.GitDetails.Action = domain.AddAction
			}
			result = append(result, rv)
		}
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

func (s *service) FormCurCommitHeaderMsg(curCommitHash string) (string, error) {
	curCommit, err := s.repo.GetCurrentCommit(curCommitHash)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("-- %s(%s)", strings.Split(curCommit.Message, "\n")[0], curCommit.Committer.Name), nil
}

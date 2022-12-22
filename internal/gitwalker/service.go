package gitwalker

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/sirupsen/logrus"
)

type Service interface {
	Patch() ([]domain.File, error)
}

type service struct {
	repo Repository
}

var _ Service = (*service)(nil)

func NewService(repository Repository) *service {
	return &service{repo: repository}
}

func (s *service) Patch() ([]domain.File, error) {
	//ee9392b1a89fff1ebd9f7148cb4130ffb79c7e6f eb54a89fa224842d5dabe89f615f834193edf0d7 eb54a89fa224842d5dabe89f615f834193edf0d7
	currentCommit, err := s.repo.GetPreviousCommit("9b7c2a074bfbf22256b3728629182fe9686c9773")
	if err != nil {
		return nil, err
	}

	headCommit, err := s.repo.GetHeadCommit()
	if err != nil {
		return nil, err
	}

	files, err := s.repo.GetFilesDiff(headCommit, currentCommit)
	if err != nil {
		return nil, err
	}

	logrus.Info(files)

	return nil, nil
}

//func (s *service) getPreviousHash(hashStr string) plumbing.Hash {
//	currentHash := plumbing.NewHash(hashStr)
//	cIter, err := r.Log(&git.LogOptions{From: currentHash})
//	if err != nil {
//		logrus.Fatalf("error while getting commit history %s", err)
//	}
//
//}

//func (s *service) GetFileChanges(path string, startHash string) ([]domain.File, error) {
//	r, err := git.PlainOpen(path)
//
//	ref, err := r.Head()
//
//	// get the commit object, pointed by ref
//	commit, err := r.CommitObject(ref.Hash())
//
//	// retrieves the commit history
//	history, err := commit.History()
//
//	// iterates over the commits and print each
//	for _, c := range history {
//		fmt.Println(c)
//	}
//}

package gitwalker

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	GetHeadCommit() (*object.Commit, error)
	GetPreviousCommit(hashStr string) (*object.Commit, error)
	GetFilesChanged(form, to *object.Commit) error
}

type repository struct {
	gitRepo *git.Repository
}

var _ Repository = (*repository)(nil)

func NewRepository(gitRepo *git.Repository) *repository {
	return &repository{gitRepo: gitRepo}
}

func (r *repository) GetHeadCommit() (*object.Commit, error) {
	ref, err := r.gitRepo.Head()
	if err != nil {
		return nil, err
	}

	rCommit, err := r.gitRepo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}
	logrus.Infof("Head commit hash %s", rCommit.Hash)
	return rCommit, nil
}

func (r *repository) GetPreviousCommit(hashStr string) (*object.Commit, error) {
	currentHash := plumbing.NewHash(hashStr)
	cIter, err := r.gitRepo.Log(&git.LogOptions{From: currentHash})
	if err != nil {
		return nil, err
	}
	defer cIter.Close()

	_, err = cIter.Next()
	if err != nil {
		return nil, err
	}

	nextCommit, err := cIter.Next()
	if err != nil {
		return nil, err
	}
	logrus.Info(nextCommit.Hash.String())
	return nextCommit, nil
}

//func (r *repository) GetPreviousCommit(hashStr string) (*object.Commit, error) {
//	var nextCommit *object.Commit
//	//	currentHash := plumbing.NewHash(hashStr)
//
//	ref, err := r.gitRepo.Head()
//	if err != nil {
//		return nil, err
//	}
//
//	//rCommit, err := r.gitRepo.CommitObject(ref.Hash())
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//var isValid object.CommitFilter = func(commit *object.Commit) bool {
//	//	//_, ok := seen[commit.Hash]
//	//
//	//	// len(commit.ParentHashes) filters out merge commits
//	//	//return len(commit.ParentHashes) < 2
//	//	return true
//	//}
//
//	cIter, err := r.gitRepo.Log(&git.LogOptions{From: ref.Hash()})
//	//cIter := object.NewFilterCommitIter(rCommit, &isValid, nil)
//
//	//cIter, err := r.gitRepo.CommitObjects()
//	if err != nil {
//		return nil, err
//	}
//	defer cIter.Close()
//
//	counter := 0
//	if err := cIter.ForEach(func(commit *object.Commit) error {
//		if counter < 5 {
//			logrus.Infof("len(commit.ParentHashes) %d", len(commit.ParentHashes))
//			logrus.Infof("commit %+v", commit)
//			commit.Files()
//		}
//
//		if counter == 2 {
//			nextCommit = commit
//		}
//		counter++
//		return nil
//	}); err != nil {
//		return nil, err
//	}
//
//	logrus.Info(nextCommit.Hash.String())
//	return nextCommit, nil
//}

func (r *repository) GetFilesChanged(from, to *object.Commit) error {
	patch, err := from.Patch(to)
	if err != nil {
		return err
	}
	for _, val := range patch.FilePatches() {
		if val != nil {
			from, to := val.Files()
			if from != nil && to != nil {
				if from.Path() != to.Path() {
					fmt.Printf("File renamed from %s to %s \n", from.Path(), to.Path())
				} else {
					fmt.Printf("File changed %s\n", from.Path())
				}
			} else if from != nil && to == nil {
				fmt.Printf("File deleted %s\n", from.Path())
			} else if to != nil && from == nil {
				fmt.Printf("File created %s\n", to.Path())
			}
		}
	}

	return nil
}

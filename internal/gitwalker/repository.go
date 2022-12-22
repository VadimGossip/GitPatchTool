package gitwalker

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	GetHeadCommit() (*object.Commit, error)
	GetPreviousCommit(hashStr string) (*object.Commit, error)
	GetFilesDiff(from, to *object.Commit) ([]domain.File, error)
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

func (r *repository) addFileChanges(nextCommit, currentCommit *object.Commit, files *[]domain.File) error {
	patch, err := currentCommit.Patch(nextCommit)
	if err != nil {
		return err
	}
	for _, val := range patch.FilePatches() {
		if val != nil {
			fromFile, toFile := val.Files()
			if fromFile != nil && toFile != nil {
				file := domain.File{
					Name:        toFile.Path(),
					InitialName: fromFile.Path(),
					Path:        toFile.Path(),
					Type:        domain.OracleFileType,
				}
				if fromFile.Path() != toFile.Path() {
					file.Action = domain.RenameAction
				} else {
					file.Action = domain.ModifyAction
				}
				*files = append(*files, file)
			} else if fromFile != nil && toFile == nil {
				file := domain.File{
					Name: fromFile.Path(),
					Path: fromFile.Path(),
					Type: domain.OracleFileType,
				}
				file.Action = domain.DeleteAction
				*files = append(*files, file)
			} else if toFile != nil && fromFile == nil {
				file := domain.File{
					Name: toFile.Path(),
					Path: toFile.Path(),
					Type: domain.OracleFileType,
				}
				file.Action = domain.AddAction
				*files = append(*files, file)
			}
		}
	}
	return nil
}

func (r *repository) GetFilesDiff(from, to *object.Commit) ([]domain.File, error) {
	files := make([]domain.File, 0)
	commitIter, err := r.gitRepo.Log(&git.LogOptions{From: from.Hash, Since: &to.Committer.When})

	if err := commitIter.ForEach(func(commit *object.Commit) error {
		if len(commit.ParentHashes) < 2 {
			err = r.addFileChanges(from, commit, &files)
			if err != nil {
				return err
			}
		}
		from = commit
		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

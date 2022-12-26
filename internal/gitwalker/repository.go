package gitwalker

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
	"path/filepath"
	"strings"
)

type Repository interface {
	GetHeadCommit() (*object.Commit, error)
	GetPreviousCommit(hashStr string) (*object.Commit, error)
	GetFilesDiff(head, till *object.Commit) ([]domain.File, error)
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
	return nextCommit, nil
}

func (r *repository) addFileChanges(nextCommit, currentCommit *object.Commit, files *[]domain.File) error {
	patch, err := currentCommit.Patch(nextCommit)
	if err != nil {
		return err
	}
	for _, val := range patch.FilePatches() {
		if val != nil {
			var fromFileName, toFileName string
			fromFile, toFile := val.Files()
			if fromFile != nil {
				fromFileName = filepath.Base(fromFile.Path())
				if filepath.Ext(fromFile.Path()) != ".sql" {
					fromFile = nil
				}
			}
			if toFile != nil {
				toFileName = filepath.Base(toFile.Path())
				if filepath.Ext(toFile.Path()) != ".sql" {
					toFile = nil
				}
			}

			if fromFile != nil && toFile != nil {
				file := domain.File{
					Name:        fromFileName,
					InitialName: toFileName,
					Path:        strings.Replace(toFile.Path(), "/", string(os.PathSeparator), -1),
					Type:        domain.OracleFileType,
				}
				if fromFile.Path() != toFile.Path() {
					file.GitAction = domain.RenameAction
				} else {
					file.GitAction = domain.ModifyAction
				}
				*files = append(*files, file)
			} else if fromFile != nil && toFile == nil {
				file := domain.File{
					Name: fromFileName,
					Path: strings.Replace(fromFile.Path(), "/", string(os.PathSeparator), -1),
					Type: domain.OracleFileType,
				}
				file.GitAction = domain.DeleteAction
				*files = append(*files, file)
			} else if toFile != nil && fromFile == nil {
				file := domain.File{
					Name: toFileName,
					Path: strings.Replace(toFile.Path(), "/", string(os.PathSeparator), -1),
					Type: domain.OracleFileType,
				}
				file.GitAction = domain.AddAction
				*files = append(*files, file)
			}
		}
	}
	return nil
}

func (r *repository) GetFilesDiff(head, till *object.Commit) ([]domain.File, error) {
	files := make([]domain.File, 0)
	commitIter, err := r.gitRepo.Log(&git.LogOptions{From: head.Hash, Since: &till.Committer.When})

	if err := commitIter.ForEach(func(commit *object.Commit) error {
		if len(commit.ParentHashes) < 2 {
			err = r.addFileChanges(head, commit, &files)
			if err != nil {
				return err
			}
		}
		head = commit
		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

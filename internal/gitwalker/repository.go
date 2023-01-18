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
	GetCurrentCommit(hashStr string) (*object.Commit, error)
	GetPreviousCommit(hashStr string) (*object.Commit, error)
	GetFilesDiff(head, till *object.Commit) ([]domain.File, error)
}

type repository struct {
	rootDir string
	gitRepo *git.Repository
}

var _ Repository = (*repository)(nil)

func NewRepository(rootDir string, gitRepo *git.Repository) *repository {
	return &repository{rootDir: rootDir, gitRepo: gitRepo}
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

func (r *repository) GetCurrentCommit(hashStr string) (*object.Commit, error) {
	currentHash := plumbing.NewHash(hashStr)
	return r.gitRepo.CommitObject(currentHash)
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
				file := domain.File{}
				if fromFile.Path() != toFile.Path() {
					file = domain.File{
						Name: toFileName,
						Path: strings.Replace(r.rootDir+toFile.Path(), "/", string(os.PathSeparator), -1),
						GitDetails: domain.GitFileDetails{
							InitialName: fromFileName,
							InitialPath: strings.Replace(r.rootDir+fromFile.Path(), "/", string(os.PathSeparator), -1),
							Comment:     strings.Split(currentCommit.Message, "\n")[0],
							Action:      domain.RenameAction,
						},
					}
				} else {
					file = domain.File{
						Name: toFileName,
						Path: strings.Replace(r.rootDir+toFile.Path(), "/", string(os.PathSeparator), -1),
						GitDetails: domain.GitFileDetails{
							Comment: strings.Split(currentCommit.Message, "\n")[0],
							Action:  domain.ModifyAction,
						},
					}
				}
				*files = append(*files, file)
			} else if fromFile != nil && toFile == nil {
				file := domain.File{
					Name: fromFileName,
					Path: strings.Replace(r.rootDir+fromFile.Path(), "/", string(os.PathSeparator), -1),
					GitDetails: domain.GitFileDetails{
						Comment: strings.Split(currentCommit.Message, "\n")[0],
						Action:  domain.DeleteAction,
					},
				}
				*files = append(*files, file)
			} else if toFile != nil && fromFile == nil {
				file := domain.File{
					Name: toFileName,
					Path: strings.Replace(r.rootDir+toFile.Path(), "/", string(os.PathSeparator), -1),
					GitDetails: domain.GitFileDetails{
						Comment: strings.Split(currentCommit.Message, "\n")[0],
						Action:  domain.AddAction,
						New:     true,
					},
				}
				*files = append(*files, file)
			}
		}
	}
	return nil
}

func (r *repository) commitSuitable(commit object.Commit) bool {
	return !strings.HasPrefix(strings.ToLower(commit.Message), "merge")
}

func (r *repository) GetFilesDiff(headCommit, fromCommit *object.Commit) ([]domain.File, error) {
	files := make([]domain.File, 0)
	orderedFiles := make([]domain.File, 0)
	commitIter, err := r.gitRepo.Log(&git.LogOptions{From: headCommit.Hash, Order: git.LogOrderCommitterTime, Since: &fromCommit.Committer.When})

	if err := commitIter.ForEach(func(commit *object.Commit) error {
		if r.commitSuitable(*headCommit) && headCommit.Hash != commit.Hash {
			if err = r.addFileChanges(headCommit, commit, &files); err != nil {
				return err
			}
		}
		headCommit = commit
		return nil
	}); err != nil {
		return nil, err
	}

	for i := len(files) - 1; i >= 0; i-- {
		orderedFiles = append(orderedFiles, files[i])
	}
	return orderedFiles, nil
}

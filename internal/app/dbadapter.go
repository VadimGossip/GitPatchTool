package app

import (
	"fmt"
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/gitwalker"
	"github.com/go-git/go-git/v5"
)

type DBAdapter struct {
	gitWalkerRepo gitwalker.Repository
	cfg           *domain.Config
}

func NewDBAdapter(cfg *domain.Config) *DBAdapter {
	return &DBAdapter{cfg: cfg}
}

func (d *DBAdapter) Connect() error {

	gitRepo, err := git.PlainOpen(d.cfg.Path.RootDir)
	if err != nil {
		return fmt.Errorf("error while opening git repo %s", err)
	}
	d.gitWalkerRepo = gitwalker.NewRepository(d.cfg.Path.RootDir, gitRepo)
	return nil
}

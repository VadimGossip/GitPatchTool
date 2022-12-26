package app

import (
	"fmt"
	"github.com/VadimGossip/gitPatchTool/internal/gitwalker"
	"github.com/go-git/go-git/v5"
)

type DBAdapter struct {
	gitWalkerRepo gitwalker.Repository
}

func NewDBAdapter() *DBAdapter {
	return &DBAdapter{}
}

func (d *DBAdapter) Connect() error {

	gitRepo, err := git.PlainOpen("e:\\WorkSpace\\TCS_Oracle\\")
	if err != nil {
		return fmt.Errorf("error while opening git repo %s", err)
	}
	d.gitWalkerRepo = gitwalker.NewRepository("e:\\WorkSpace\\TCS_Oracle\\", gitRepo)
	return nil
}

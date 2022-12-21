package patcher

import (
	"fmt"
	"github.com/VadimGossip/gitPatchTool/internal/filewalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
)

type Service interface {
	CreatePatch() error
	FixObjectFile() error
}

type service struct {
	fileWalker filewalker.Service
	extractor  extractor.Service
}

var _ Service = (*service)(nil)

func NewService(fileWalker filewalker.Service, extractor extractor.Service) *service {
	return &service{fileWalker: fileWalker, extractor: extractor}
}

func (s *service) CreatePatch() error {

	// We instantiate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen("e:\\WorkSpace\\TCS_Oracle\\")
	if err != nil {
		logrus.Fatalf("error while opening git repo %s", err)
	}
	//git.PlainOpenWithOptions()

	// ... retrieving the HEAD reference
	ref, err := r.Head()
	if err != nil {
		logrus.Fatalf("error while getting head %s", err)
	}

	commitId := "95310e5bab4d33c1f401db0deb4e97e5108a356a"
	ss := plumbing.NewHash(commitId)

	h, _ := r.CommitObject(ss)
	refc, _ := r.CommitObject(ref.Hash())

	time := h.Committer.When
	logrus.Info(h.Committer.When)

	xxx, _ := h.Patch(refc)
	///fmt.Println(xxx)
	//fmt.Println(xxx.Stats())
	for _, val := range xxx.FilePatches() {
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

	//tIter,err := r.TreeObjects()
	//if err != nil {
	//	logrus.Fatalf("error while getting tree history %s", err)
	//}
	//err = tIter.ForEach(func(c *object.Tree) error {
	//	c.
	//	//fmt.Println(c)
	//	return nil
	//}))

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &time})
	if err != nil {
		logrus.Fatalf("error while getting commit history %s", err)
	}

	defer cIter.Close()
	// ... just iterates over the commits
	var cCount int
	err = cIter.ForEach(func(cx *object.Commit) error {
		//stats, _ := c.Stats()
		//for _, st := range stats {
		//	fmt.Println(st.Name)
		//}
		//err = iter.ForEach(func(f *object.File) error {
		//	fmt.Printf("file %s\n", f.Name)
		//	return nil
		//})
		//patch, _ := c.Patch(h)
		//logrus.Println(patch)
		//err = tree.Files().ForEach(func(f *object.File) error {
		//	fmt.Printf("file %s\n", f.Name)
		//	return nil
		//})
		cCount++
		return nil
	})
	if err != nil {
		logrus.Fatalf("error while iterating commits git repo %s", err)
	}

	logrus.Println(cCount)

	//files, err := s.fileWalker.Walk("e:\\WorkSpace\\TCS_Oracle\\", domain.OracleFileType)
	//if err != nil {
	//	logrus.Fatalf("Fail to collect files %s", err)
	//}
	//
	//oraObjects := s.extractor.ExtractOracleObjects(files)
	//
	//for _, obj := range oraObjects {
	//	if len(obj.Errors) > 0 {
	//		logrus.Info(obj)
	//	}
	//}

	return nil
}

func (s *service) FixObjectFile() error {
	return nil
}

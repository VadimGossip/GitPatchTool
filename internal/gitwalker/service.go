package gitwalker

type Service interface {
}

type service struct {
}

var _ Service = (*service)(nil)

func NewService() *service {
	return &service{}
}

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

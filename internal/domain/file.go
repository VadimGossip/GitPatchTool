package domain

const (
	OracleFileType int = 1
)

const (
	AddAction int = iota
	DeleteAction
	ModifyAction
	RenameAction
)

type File struct {
	Name             string
	ShortPath        string
	Path             string
	InitialName      string
	InitialShortPath string
	InitialPath      string
	Type             int
	GitAction        int
}

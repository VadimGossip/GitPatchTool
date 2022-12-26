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
	Name        string
	InitialName string
	ShortPath   string
	Path        string
	Type        int
	GitAction   int
}

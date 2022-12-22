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
	Path        string
	Type        int
	Action      int
}

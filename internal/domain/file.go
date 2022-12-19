package domain

const (
	OracleFileType int = 1
)

type File struct {
	Name string
	Path string
	Type int
}

type OracleObject struct {
	Name string
	Path string
	Type int
}

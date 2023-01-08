package domain

const (
	PatchMode string = "patch"
)

type Path struct {
	RootDir    string
	InstallDir string
}

type Config struct {
	Mode     string
	CommitId string
	Path     Path
}

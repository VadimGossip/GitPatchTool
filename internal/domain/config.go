package domain

const (
	PatchMode string = "patch"
)

type Path struct {
	RootDir    string
	InstallDir string
}

type DictionariesConfig struct {
	ServerSchema      map[string]ServerSchema
	InstallFilename   map[ServerSchema]string
	MigrationFilename map[ServerSchema]string
	ServerSchemaAlias map[ServerSchema]string
}

type Config struct {
	Mode         string
	CommitId     string
	Path         Path
	Dictionaries DictionariesConfig
}

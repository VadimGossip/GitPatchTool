package domain

const (
	OracleFileType int = 1
)

const (
	AddAction    int = 1
	DeleteAction int = 2
	ModifyAction int = 3
	RenameAction int = 4
)

const (
	Ordinary int = 1
	ErrorLog int = 2
	Warning  int = 3
)

const (
	ErrorLogFileName   string = "error_log.txt"
	WarningLogFileName string = "warning_log.txt"
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

type InstallFile struct {
	Path      string
	FileLines []string
	Type      int
}

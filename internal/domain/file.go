package domain

const (
	AddAction    int = 1
	DeleteAction int = 2
	ModifyAction int = 3
	RenameAction int = 4
)

const (
	Data int = iota
	Install
	ErrorLog
	WarningLog
)

const (
	ErrorLogFileName   string = "error_log.txt"
	WarningLogFileName string = "warning_log.txt"
)

// GitFileDetails add and handle initialAction to handle add -> modify -> delete sequence
type GitFileDetails struct {
	InitialName string
	InitialPath string
	Comment     string
	Action      int
	New         bool
}

type File struct {
	Name       string
	Path       string
	FileLines  []string
	GitDetails GitFileDetails
}

type OracleFile struct {
	OracleDataType int
	FileDetails    File
}

var ActionNameDict = map[int]string{
	AddAction:    "Add",
	DeleteAction: "Delete",
	ModifyAction: "Modify",
	RenameAction: "Rename",
}

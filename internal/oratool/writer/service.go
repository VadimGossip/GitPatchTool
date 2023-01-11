package writer

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
)

type Service interface {
	CreateInstallLines(installDir string, oracleObjects []domain.OracleObject) []domain.InstallFile
}

type service struct {
}

var _ Service = (*service)(nil)

func NewService() *service {
	return &service{}
}

func (s *service) formErrorLines(installDir string, oracleObjects []domain.OracleObject) domain.InstallFile {
	objErrors := make(map[string][]domain.OracleObject)
	errorLines := make([]string, 0)

	for _, val := range oracleObjects {
		for _, errMsg := range val.Errors {
			objErrors[errMsg] = append(objErrors[errMsg], val)
		}
	}
	for _, errMsg := range []string{domain.UnknownObjectType.Error(), domain.SchemaNotFound.Error(), domain.FileNotExists.Error()} {
		for idx, val := range objErrors[errMsg] {
			if idx == 0 {
				errorLines = append(errorLines, errMsg)
			}
			errorLines = append(errorLines, val.File.Path)
		}
	}

	return domain.InstallFile{
		Path:      installDir + domain.ErrorLogFileName,
		FileLines: errorLines,
		Type:      domain.ErrorLog,
	}
}

func (s *service) CreateInstallLines(installDir string, oracleObjects []domain.OracleObject) []domain.InstallFile {
	resultFiles := make([]domain.InstallFile, 0)
	objWErrors := make([]domain.OracleObject, 0)

	for _, obj := range oracleObjects {
		if len(obj.Errors) != 0 {
			objWErrors = append(objWErrors, obj)
		}
	}

	if len(objWErrors) > 1 {
		resultFiles = append(resultFiles, s.formErrorLines(installDir, objWErrors))
	}

	return resultFiles
}

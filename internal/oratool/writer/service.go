package writer

import (
	"fmt"
	"github.com/VadimGossip/gitPatchTool/internal/domain"
)

type Service interface {
	CreateInstallLines(oracleObjects []domain.OracleObject) []domain.InstallFile
}

type service struct {
}

var _ Service = (*service)(nil)

func NewService() *service {
	return &service{}
}

func (s *service) formErrorLines(oracleObjects []domain.OracleObject) domain.InstallFile {
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
		fmt.Println()
	}

	return domain.InstallFile{
		Path:      "",
		FileLines: errorLines,
		Type:      domain.ErrorLog,
	}
}

func (s *service) CreateInstallLines(oracleObjects []domain.OracleObject) []domain.InstallFile {
	resultFiles := make([]domain.InstallFile, 0)
	objWErrors := make([]domain.OracleObject, 0)

	for _, obj := range oracleObjects {
		if len(obj.Errors) != 0 {
			objWErrors = append(objWErrors, obj)
		}
	}

	if len(objWErrors) > 1 {
		resultFiles = append(resultFiles, s.formErrorLines(objWErrors))
	}

	return resultFiles
}

package writer

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
)

type Service interface {
	CreateInstallLines(installDir string, oracleObjects []domain.OracleObject) []domain.OracleFile
}

type service struct {
}

var _ Service = (*service)(nil)

func NewService() *service {
	return &service{}
}

func (s *service) formErrorLines(installDir string, oracleObjects []domain.OracleObject) domain.OracleFile {
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

	return domain.OracleFile{
		OracleDataType: domain.ErrorLog,
		FileDetails: domain.File{
			Name:      domain.ErrorLogFileName,
			Path:      installDir + domain.ErrorLogFileName,
			FileLines: errorLines,
		},
	}
}

//func (s *service) getInstallFileName(serverSchema []domain.ServerSchema) string {
//   return domain.GetInstallFileNameByServerSchema(serverSchema)
//
//}

//func (s *service) sortOracleObject(oracleObjects []domain.OracleObject) {
//	sort.SliceStable(oracleObjects, func(i, j int) bool {
//		if (oracleObjects[i].ObjectType == domain.OracleScriptsBeforeType &&
//			oracleObjects[j].ObjectType != domain.OracleScriptsBeforeType) ||
//			(oracleObjects[i].ObjectType != domain.OracleScriptsBeforeType && oracleObjects[i].ObjectType != domain.OracleScriptsAfterType &&
//				oracleObjects[j].ObjectType == domain.OracleScriptsAfterType){
//			return true
//		}
//
//
//		if oracleObjects[i].ObjectType == oracleObjects[j].ObjectType &&
//
//
//
//		return suppliers[i].EUM != nil && (suppliers[j].EUM == nil || *suppliers[i].EUM > *suppliers[j].EUM)
//	})
//}

//func (s *service) formInstallLines(installDir string, oracleObjects []domain.OracleObject) []domain.InstallFile {
//	objInstall := make(map[string][]domain.OracleObject)
//	//InstallLines := make([]string, 0)
//
//	for _, val := range oracleObjects {
//		for _, serverSchema := range val.ServerSchemaList {
//			installFileName := domain.GetInstallFileNameByServerSchema(serverSchema)
//			if installFileName != "" {
//				objInstall[installFileName] = append(objInstall[installFileName], val)
//			}
//		}
//	}
//	//for _, errMsg := range []string{domain.UnknownObjectType.Error(), domain.SchemaNotFound.Error(), domain.FileNotExists.Error()} {
//	//	for idx, val := range objErrors[errMsg] {
//	//		if idx == 0 {
//	//			errorLines = append(errorLines, errMsg)
//	//		}
//	//		errorLines = append(errorLines, val.File.Path)
//	//	}
//	//}
//	//
//	//return domain.InstallFile{
//	//	Path:      installDir + domain.ErrorLogFileName,
//	//	FileLines: errorLines,
//	//	Type:      domain.ErrorLog,
//	//}
//}

func (s *service) CreateInstallLines(installDir string, oracleObjects []domain.OracleObject) []domain.OracleFile {
	resultFiles := make([]domain.OracleFile, 0)
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

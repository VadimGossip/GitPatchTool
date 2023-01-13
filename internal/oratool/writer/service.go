package writer

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"sort"
	"strings"
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
			errorLines = append(errorLines, val.File.FileDetails.Path)
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

////needToReworkIn Future for schema install
//func (s *service) getInstallWeight(objectType int) int {
//	sortMap := map[int]int{
//		domain.
//	}
//

//type key struct {
//	objectType int
//	action     int
//}
//sortMap := map[key]int{
//	{
//		objectType: domain.OracleScriptsBeforeType, action: domain.AddAction} : 1,
//}
//
//vtbs_clogs delete(
//	)

//allowed_object_type_sort_mask = {"tablespaces"      : 0
//	,"directories"      : 1
//	,"dblinks"          : 2
//	,"users"            : 3
//	,"synonyms"         : 4
//	,"scripts_before"   : 5
//	,"contexts"         : 6
//	,"sequences"        : 7
//	,"types"            : 8
//	,"tables"           : 9
//	,"mlogs"            : 10
//	,"mviews"           : 11
//	,"types"            : 12
//	,"packages"         : 13
//	,"views"            : 14
//	,"triggers"         : 15
//	,"vtbs_tasks"       : 16
//	,"rows"             : 17
//	,"roles"            : 18
//	,"functions"        : 19
//	,"vtbs_clogs"       : 20
//	,"scripts_after"    : 21
//	,"scripts_migration": 22}

//OracleTablespaceType = iota
//OracleDirectoryType
//OracleDbLinkType
//OracleUserType
//OracleSynonymType
//OracleContextType
//OracleSequencesType
//OracleTypeType
//OracleTableType
//OracleMLogType
//OracleMViewType
//OraclePackageType
//OracleViewType
//OracleTriggerType
//OracleVTaskType
//OracleRowType
//OracleRoleType
//OracleFunctionType
//OracleVClogType
//OracleTableFKType
//OracleScriptsBeforeType
//OracleScriptsAfterType
//OracleScriptsMigrationType
//

//object_type_skip_set = {'install'
//	,'tables'
//	,'rows'
//	,'roles'
//	,'users'
//	,'dblinks'
//	,'tablespaces'}

//	return 0
//}

func (s *service) getPackageWeight(oracleObject domain.OracleObject) int {
	if strings.HasSuffix(oracleObject.ObjectName, "read") {
		return 0
	}
	if strings.HasSuffix(oracleObject.ObjectName, "digests") {
		return 1
	}
	if strings.HasSuffix(oracleObject.ObjectName, "utils") {
		return 2
	}
	if strings.HasSuffix(oracleObject.ObjectName, "ri") {
		return 3
	}
	if strings.HasSuffix(oracleObject.ObjectName, "ui") {
		return 4
	}
	return 5
}
func (s *service) getTypeWeight(oracleObject domain.OracleObject) int {
	return len(oracleObject.ObjectName)
}

func (s *service) sortOracleObject(oracleObjects []domain.OracleObject) {
	sort.SliceStable(oracleObjects, func(i, j int) bool {
		//if (oracleObjects[i].ObjectType == domain.OracleScriptsBeforeType && oracleObjects[j].ObjectType != domain.OracleScriptsBeforeType) ||
		//	(oracleObjects[i].ObjectType != domain.OracleScriptsBeforeType && oracleObjects[j].ObjectType == domain.OracleScriptsBeforeType) {
		//	return true
		//}

		if oracleObjects[i].ObjectType == oracleObjects[j].ObjectType && oracleObjects[i].ObjectType == domain.OraclePackageType {
			return s.getPackageWeight(oracleObjects[i]) < s.getPackageWeight(oracleObjects[j])
		}

		if oracleObjects[i].ObjectType == oracleObjects[j].ObjectType && oracleObjects[i].ObjectType == domain.OracleTypeType {
			return s.getPackageWeight(oracleObjects[i]) < s.getPackageWeight(oracleObjects[j])
		}

		return oracleObjects[i].ObjectType < oracleObjects[j].ObjectType
	})
}

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

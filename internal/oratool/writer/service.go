package writer

import (
	"fmt"
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"sort"
	"strings"
)

type Service interface {
	CreateInstallFiles(installDir string, oracleObjects []domain.OracleObject) []domain.OracleFile
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

func (s *service) sortOracleObjects(oracleObjects []domain.OracleObject) {
	sort.SliceStable(oracleObjects, func(i, j int) bool {
		if (oracleObjects[i].ObjectType == domain.OracleScriptsBeforeType && oracleObjects[j].ObjectType != domain.OracleScriptsBeforeType) ||
			(oracleObjects[i].ObjectType != domain.OracleScriptsBeforeType && oracleObjects[j].ObjectType == domain.OracleScriptsBeforeType) ||
			(oracleObjects[i].ObjectType != domain.OracleScriptsMigrationType && oracleObjects[j].ObjectType == domain.OracleScriptsMigrationType) {
			return true
		}

		if oracleObjects[i].EpicModuleName < oracleObjects[j].EpicModuleName ||
			oracleObjects[i].EpicModuleName == oracleObjects[j].EpicModuleName && oracleObjects[i].ModuleName < oracleObjects[j].ModuleName {
			return true
		}

		if oracleObjects[i].ObjectType == oracleObjects[j].ObjectType && oracleObjects[i].ModuleName == oracleObjects[j].ModuleName &&
			oracleObjects[i].ObjectType == domain.OraclePackageType {
			return s.getPackageWeight(oracleObjects[i]) < s.getPackageWeight(oracleObjects[j])
		}

		if oracleObjects[i].ObjectType == oracleObjects[j].ObjectType && oracleObjects[i].ModuleName == oracleObjects[j].ModuleName &&
			oracleObjects[i].ObjectType == domain.OracleTypeType {
			return s.getPackageWeight(oracleObjects[i]) < s.getPackageWeight(oracleObjects[j])
		}

		return oracleObjects[i].ObjectType < oracleObjects[j].ObjectType
	})
}

//func (s *service) formObjInstallLines(obj domain.OracleObject)

func (s *service) formModuleHeader(obj domain.OracleObject) string {
	if obj.ObjectType == domain.OracleScriptsBeforeType {
		return fmt.Sprintf("-------------------------%s/%s-------------------------", obj.EpicModuleName, "scripts_before")
	} else if obj.ObjectType == domain.OracleScriptsAfterType {
		return fmt.Sprintf("-------------------------%s/%s-------------------------", obj.EpicModuleName, "scripts_after")
	} else if obj.ObjectType == domain.OracleScriptsMigrationType {
		return fmt.Sprintf("-------------------------%s/%s-------------------------", obj.EpicModuleName, "scripts_migration")
	} else {
		return fmt.Sprintf("-------------------------%s/%s-------------------------", obj.EpicModuleName, obj.ModuleName)
	}
}

func (s *service) formObjectLines(obj domain.OracleObject) []string {
	return []string{
		fmt.Sprintf("prompt %s", strings.Replace(obj.File.FileDetails.Path, "", "", -1)),
		fmt.Sprintf("%s", strings.Replace(obj.File.FileDetails.Path, "@ ../..", "", -1)),
	}
}

func (s *service) formInstallFiles(installDir string, oracleObjects []domain.OracleObject) []domain.OracleFile {
	objInstall := make(map[string][]domain.OracleObject)
	installFileLines := make(map[string][]string)
	//InstallLines := make([]string, 0)

	for _, obj := range oracleObjects {
		for _, serverSchema := range obj.ServerSchemaList {
			installFileName := domain.GetInstallFileNameByServerSchema(serverSchema)
			if installFileName != "" {
				objInstall[installFileName] = append(objInstall[installFileName], obj)
			}
		}
	}

	for key, objI := range objInstall {
		s.sortOracleObjects(objI)
		//check if file not exists add file header
		//add commit header
		for idx := range objI {
			if idx == 0 || idx > 0 && s.formModuleHeader(objI[idx-1]) != s.formModuleHeader(objI[idx]) {
				installFileLines[key] = append(installFileLines[key], s.formModuleHeader(objI[idx]))
			}

			if idx == 0 || idx > 0 && domain.GetDirNameByOracleObjectType(objI[idx-1].ObjectType) != domain.GetDirNameByOracleObjectType(objI[idx].ObjectType) {
				installFileLines[key] = append(installFileLines[key], domain.GetDirNameByOracleObjectType(objI[idx].ObjectType))
			}

			installFileLines[key] = append(installFileLines[key], s.formObjectLines(objI[idx])...)

			//} else if idx > 0 && s.formModuleHeader(objI[idx-1]) != s.formModuleHeader(objI[idx]) {
			//
			//	if objI[idx-1].ModuleName != objI[idx-1].ModuleName {
			//		installFiles[key] = append(installFiles[key], "")
			//	}
			//}
		}
		//add footer
	}

	for key, ifls := range installFileLines {
		for _, ifl := range ifls {
			fmt.Printf("filename %s fileline %s\n", key, ifl)
		}
	}

	//return domain.InstallFile{
	//	Path:      installDir + domain.ErrorLogFileName,
	//	FileLines: errorLines,
	//	Type:      domain.ErrorLog,
	//}
	return nil
}

func (s *service) CreateInstallFiles(installDir string, oracleObjects []domain.OracleObject) []domain.OracleFile {
	resultFiles := make([]domain.OracleFile, 0)
	objWErrors := make([]domain.OracleObject, 0)
	objInstall := make([]domain.OracleObject, 0)

	for _, obj := range oracleObjects {
		if len(obj.Errors) != 0 {
			objWErrors = append(objWErrors, obj)
		} else if obj.File.OracleDataType == domain.Data {
			objInstall = append(objInstall, obj)
		}

	}

	if len(objWErrors) > 1 {
		resultFiles = append(resultFiles, s.formErrorLines(installDir, objWErrors))
	}

	if len(objInstall) > 1 {
		resultFiles = append(resultFiles, s.formInstallFiles(installDir, objInstall)...)
	}

	return resultFiles
}

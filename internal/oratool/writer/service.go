package writer

import (
	"fmt"
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/file"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Service interface {
	CreateInstallFiles(rootDir, installDir, commitMsg string, oracleObjects []domain.OracleObject) error
}

type service struct {
	file file.Service
}

var _ Service = (*service)(nil)

func NewService(file file.Service) *service {
	return &service{file: file}
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

func (s *service) createInstallFileHeader(schemaStrItem, filename string) []string {
	return []string{
		fmt.Sprintf("-- Schema: %s\n", schemaStrItem),
		fmt.Sprintf("prompt install %s\n", filename),
		fmt.Sprintf("set define off\n"),
		fmt.Sprintf("spool %s.log append\n", filename[:len(filename)-len(filepath.Ext(filename))]),
	}
}

func (s *service) removeLineWithPrefix(lines []string, prefix string) []string {
	result := make([]string, 0)
	for _, line := range lines {
		if !strings.HasPrefix(line, prefix) {
			result = append(result, line)
		}
	}
	return result
}

func (s *service) createInstallFileFooter() string {
	return "spool off"
}

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

func (s *service) formObjectLines(rootDir string, obj domain.OracleObject) []string {
	return []string{
		fmt.Sprintf("prompt %s", strings.Replace(strings.Replace(obj.File.FileDetails.Path, rootDir, "", -1), string(os.PathSeparator), "/", -1)),
		fmt.Sprintf("%s", strings.Replace(strings.Replace(obj.File.FileDetails.Path, rootDir, "@ ../../", -1), string(os.PathSeparator), "/", -1)),
	}
}
func (s *service) formErrorFile(installDir string, oracleObjects []domain.OracleObject) domain.OracleFile {
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

func (s *service) formInstallFiles(rootDir, installDir, commitMsg string, oracleObjects []domain.OracleObject) []domain.OracleFile {
	type oiKey struct {
		filename      string
		schemaStrItem string
	}
	objInstall := make(map[oiKey][]domain.OracleObject)

	installFileLines := make(map[string][]string)
	installFiles := make([]domain.OracleFile, 0)
	fileLines := make([]string, 0)

	for _, obj := range oracleObjects {
		for _, serverSchema := range obj.ServerSchemaList {
			installFileName := domain.ServerSchemaInstallFilenameDict[serverSchema]
			schemaStrItem := domain.ServerSchemaSchemaStrItemDict[serverSchema]
			if installFileName != "" {
				objInstall[oiKey{filename: installFileName, schemaStrItem: schemaStrItem}] = append(objInstall[oiKey{filename: installFileName, schemaStrItem: schemaStrItem}], obj)
			}
		}
	}

	var curModuleH, prevModuleH, curObjectTypeH, prevObjectTypeH string
	for key, objI := range objInstall {
		s.sortOracleObjects(objI)
		addToFile := s.file.CheckFileExists(installDir + key.filename)
		footer := s.createInstallFileFooter()
		if addToFile {
			fileLines, _ = s.file.ReadFileLines(installDir + key.filename)
			installFileLines[key.filename] = append(installFileLines[key.filename], s.removeLineWithPrefix(fileLines, footer)...)
		} else {
			installFileLines[key.filename] = append(installFileLines[key.filename], s.createInstallFileHeader(key.filename, key.schemaStrItem)...)
		}
		installFileLines[key.filename] = append(installFileLines[key.filename], commitMsg)
		for idx := range objI {
			curModuleH = s.formModuleHeader(objI[idx])
			curObjectTypeH = domain.OracleObjectTypeDirDict[objI[idx].ObjectType]
			if idx > 0 {
				prevModuleH = s.formModuleHeader(objI[idx-1])
				prevObjectTypeH = domain.OracleObjectTypeDirDict[objI[idx-1].ObjectType]
			}

			if curModuleH != prevModuleH {
				installFileLines[key.filename] = append(installFileLines[key.filename], "")
				installFileLines[key.filename] = append(installFileLines[key.filename], curModuleH)
			}

			if curModuleH != prevModuleH || curObjectTypeH != prevObjectTypeH {
				installFileLines[key.filename] = append(installFileLines[key.filename], "")
				installFileLines[key.filename] = append(installFileLines[key.filename], curObjectTypeH)
			}

			installFileLines[key.filename] = append(installFileLines[key.filename], s.formObjectLines(rootDir, objI[idx])...)
		}
		installFileLines[key.filename] = append(installFileLines[key.filename], "")
		installFileLines[key.filename] = append(installFileLines[key.filename], footer)
	}

	for key, lines := range installFileLines {
		installFiles = append(installFiles, domain.OracleFile{OracleDataType: domain.Install,
			FileDetails: domain.File{
				Name:      key,
				Path:      installDir + key,
				FileLines: lines,
			},
		})
		for _, ifl := range lines {
			fmt.Printf("filename %s fileline %s\n", key, ifl)
		}
	}
	return nil
}

func (s *service) CreateInstallFiles(rootDir, installDir, commitMsg string, oracleObjects []domain.OracleObject) error {
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
		errFile := s.formErrorFile(installDir, objWErrors)
		if err := s.file.DeleteFile(errFile.FileDetails.Path); err != nil {
			return err
		}
		if err := s.file.CreateFile(errFile.FileDetails.Path, errFile.FileDetails.FileLines); err != nil {
			return err
		}
	}

	if len(objInstall) > 1 {
		for _, installFile := range s.formInstallFiles(rootDir, installDir, commitMsg, objInstall) {
			if err := s.file.CreateFile(installFile.FileDetails.Path, installFile.FileDetails.FileLines); err != nil {
				return err
			}
		}
	}

	return nil
}

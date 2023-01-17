package writer

import (
	"fmt"
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/file"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

type typeAction struct {
	oracleObjectType int
	action           int
}

var manualOperations = map[typeAction]struct{}{
	//Modify Action
	{domain.OracleTablespaceType, domain.ModifyAction}:       {},
	{domain.OracleDirectoryType, domain.ModifyAction}:        {},
	{domain.OracleDbLinkType, domain.ModifyAction}:           {},
	{domain.OracleUserType, domain.ModifyAction}:             {},
	{domain.OracleVTaskType, domain.ModifyAction}:            {},
	{domain.OracleRowType, domain.ModifyAction}:              {},
	{domain.OracleRoleType, domain.ModifyAction}:             {},
	{domain.OracleTableFKType, domain.ModifyAction}:          {},
	{domain.OracleScriptsBeforeType, domain.ModifyAction}:    {},
	{domain.OracleScriptsAfterType, domain.ModifyAction}:     {},
	{domain.OracleScriptsMigrationType, domain.ModifyAction}: {},

	//Delete Action
	{domain.OracleTablespaceType, domain.DeleteAction}:       {},
	{domain.OracleDirectoryType, domain.DeleteAction}:        {},
	{domain.OracleDbLinkType, domain.DeleteAction}:           {},
	{domain.OracleUserType, domain.DeleteAction}:             {},
	{domain.OracleSynonymType, domain.DeleteAction}:          {},
	{domain.OracleContextType, domain.DeleteAction}:          {},
	{domain.OracleSequencesType, domain.DeleteAction}:        {},
	{domain.OracleTypeType, domain.DeleteAction}:             {},
	{domain.OracleTableType, domain.DeleteAction}:            {},
	{domain.OracleMLogType, domain.DeleteAction}:             {},
	{domain.OracleMViewType, domain.DeleteAction}:            {},
	{domain.OraclePackageType, domain.DeleteAction}:          {},
	{domain.OracleViewType, domain.DeleteAction}:             {},
	{domain.OracleTriggerType, domain.DeleteAction}:          {},
	{domain.OracleVTaskType, domain.DeleteAction}:            {},
	{domain.OracleRowType, domain.DeleteAction}:              {},
	{domain.OracleRoleType, domain.DeleteAction}:             {},
	{domain.OracleFunctionType, domain.DeleteAction}:         {},
	{domain.OracleVClogType, domain.DeleteAction}:            {},
	{domain.OracleTableFKType, domain.DeleteAction}:          {},
	{domain.OracleScriptsBeforeType, domain.DeleteAction}:    {},
	{domain.OracleScriptsAfterType, domain.DeleteAction}:     {},
	{domain.OracleScriptsMigrationType, domain.DeleteAction}: {},

	//RenameAction
	{domain.OracleTablespaceType, domain.RenameAction}:       {},
	{domain.OracleDirectoryType, domain.RenameAction}:        {},
	{domain.OracleDbLinkType, domain.RenameAction}:           {},
	{domain.OracleUserType, domain.RenameAction}:             {},
	{domain.OracleSynonymType, domain.RenameAction}:          {},
	{domain.OracleContextType, domain.RenameAction}:          {},
	{domain.OracleSequencesType, domain.RenameAction}:        {},
	{domain.OracleTypeType, domain.RenameAction}:             {},
	{domain.OracleTableType, domain.RenameAction}:            {},
	{domain.OracleMLogType, domain.RenameAction}:             {},
	{domain.OracleMViewType, domain.RenameAction}:            {},
	{domain.OraclePackageType, domain.RenameAction}:          {},
	{domain.OracleViewType, domain.RenameAction}:             {},
	{domain.OracleTriggerType, domain.RenameAction}:          {},
	{domain.OracleVTaskType, domain.RenameAction}:            {},
	{domain.OracleRowType, domain.RenameAction}:              {},
	{domain.OracleRoleType, domain.RenameAction}:             {},
	{domain.OracleFunctionType, domain.RenameAction}:         {},
	{domain.OracleVClogType, domain.RenameAction}:            {},
	{domain.OracleTableFKType, domain.RenameAction}:          {},
	{domain.OracleScriptsBeforeType, domain.RenameAction}:    {},
	{domain.OracleScriptsAfterType, domain.RenameAction}:     {},
	{domain.OracleScriptsMigrationType, domain.RenameAction}: {},
}

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

func (s *service) createInstallFileHeader(filename, schemaStrItem string) []string {
	return []string{
		fmt.Sprintf("-- Schema: %s", schemaStrItem),
		fmt.Sprintf("prompt install %s", filename),
		fmt.Sprintf("set define off"),
		fmt.Sprintf("spool %s.log append", filename[:len(filename)-len(filepath.Ext(filename))]),
		fmt.Sprintf(""),
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

func (s *service) formObjectTypeHeader(oracleObjectType int) string {
	return fmt.Sprintf("prompt %s", domain.OracleObjectTypeDirDict[oracleObjectType])
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

	for _, obj := range oracleObjects {
		for _, errMsg := range obj.Errors {
			objErrors[errMsg] = append(objErrors[errMsg], obj)
		}
	}
	for _, errMsg := range []string{domain.UnknownObjectType.Error(), domain.SchemaNotFound.Error(), domain.FileNotExists.Error()} {
		for idx, e := range objErrors[errMsg] {
			if idx == 0 {
				errorLines = append(errorLines, errMsg)
			}
			errorLines = append(errorLines, e.File.FileDetails.Path)
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

func (s *service) formWarningFile(installDir string, oracleObjects []domain.OracleObject) domain.OracleFile {
	objWarnings := make(map[int][]domain.OracleObject)
	warningLines := make([]string, 0)

	for _, obj := range oracleObjects {
		objWarnings[obj.File.FileDetails.GitDetails.Action] = append(objWarnings[obj.File.FileDetails.GitDetails.Action], obj)
	}
	for _, action := range []int{domain.AddAction, domain.ModifyAction, domain.DeleteAction, domain.RenameAction} {
		for idx, w := range objWarnings[action] {
			if idx == 0 {
				warningLines = append(warningLines, domain.ActionNameDict[action])
				warningLines = append(warningLines, "")
			}
			warningLines = append(warningLines, w.File.FileDetails.Path)
		}
	}

	return domain.OracleFile{
		OracleDataType: domain.WarningLog,
		FileDetails: domain.File{
			Name:      domain.WarningLogFileName,
			Path:      installDir + domain.WarningLogFileName,
			FileLines: warningLines,
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
			schemaStrItem := cases.Title(language.English, cases.Compact).String(domain.ServerSchemaSchemaStrItemDict[serverSchema])
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
			curObjectTypeH = s.formObjectTypeHeader(objI[idx].ObjectType)
			if idx > 0 {
				prevModuleH = s.formModuleHeader(objI[idx-1])
				prevObjectTypeH = s.formObjectTypeHeader(objI[idx-1].ObjectType)
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
	}
	return installFiles
}

func (s *service) removeSessionFiles(installDir string) error {
	if err := s.file.DeleteFile(installDir + domain.ErrorLogFileName); err != nil {
		return err
	}
	if err := s.file.DeleteFile(installDir + domain.WarningLogFileName); err != nil {
		return err
	}
	return nil
}

func (s *service) CreateInstallFiles(rootDir, installDir, commitMsg string, oracleObjects []domain.OracleObject) error {
	objErrors := make([]domain.OracleObject, 0)
	objWarnings := make([]domain.OracleObject, 0)
	objInstall := make([]domain.OracleObject, 0)

	if err := s.removeSessionFiles(installDir); err != nil {
		return err
	}

	for _, obj := range oracleObjects {
		if obj.File.OracleDataType == domain.Data {
			if len(obj.Errors) != 0 {
				objErrors = append(objErrors, obj)
			} else {
				if _, ok := manualOperations[typeAction{obj.ObjectType, obj.File.FileDetails.GitDetails.Action}]; !ok {
					objInstall = append(objInstall, obj)
				} else {
					objWarnings = append(objWarnings, obj)
				}
			}
		}
	}

	if len(objErrors) > 0 {
		errFile := s.formErrorFile(installDir, objErrors)
		if err := s.file.CreateFile(errFile.FileDetails.Path, errFile.FileDetails.FileLines); err != nil {
			return err
		}
	}

	if len(objWarnings) > 0 {
		warningFile := s.formWarningFile(installDir, objWarnings)
		if err := s.file.CreateFile(warningFile.FileDetails.Path, warningFile.FileDetails.FileLines); err != nil {
			return err
		}
	}

	if len(objInstall) > 0 {
		for _, installFile := range s.formInstallFiles(rootDir, installDir, commitMsg, objInstall) {
			if err := s.file.CreateFile(installFile.FileDetails.Path, installFile.FileDetails.FileLines); err != nil {
				return err
			}
		}
	}

	return nil
}

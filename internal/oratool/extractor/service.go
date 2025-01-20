package extractor

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/file"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/schema"
)

type Service interface {
	CreateOracleObjects(rootDir, installDir string, files []domain.File) []domain.OracleObject
}

type service struct {
	file   file.Service
	schema schema.Service
}

var _ Service = (*service)(nil)

func NewService(file file.Service, schema schema.Service) *service {
	return &service{
		file:   file,
		schema: schema,
	}
}

func (s *service) getObjectTypeFromDir(objectTypeDir, objectName string) (int, error) {
	if val, ok := domain.DirOracleObjectTypeDict[objectTypeDir]; ok {
		if len(strings.Split(objectName, ".")) > 1 && val == domain.OracleTableType {
			return domain.OracleTriggerType, nil
		}
		return val, nil
	}

	return 0, domain.UnknownObjectType
}

func (s *service) writeError(obj *domain.OracleObject, err error) {
	if err != nil {
		obj.Errors = append(obj.Errors, err.Error())
	}
}

func (s *service) fileSuitable(rootDir, installDir, path string) bool {
	parts := strings.Split(strings.Replace(path, rootDir, "", -1), string(os.PathSeparator))
	notInInstall := strings.ToLower(parts[0]) == "install" && (!strings.HasPrefix(path, installDir))
	inInstallButNotScript := len(parts) >= 2 && strings.HasPrefix(path, installDir) && !(domain.DirOracleObjectTypeDict[parts[len(parts)-2]] == domain.OracleScriptsMigrationType ||
		domain.DirOracleObjectTypeDict[parts[len(parts)-2]] == domain.OracleScriptsBeforeType ||
		domain.DirOracleObjectTypeDict[parts[len(parts)-2]] == domain.OracleScriptsAfterType)
	return !(notInInstall || inInstallButNotScript)
}

func (s *service) resolveAdditionalPathInfo(oracleObj *domain.OracleObject) {
	var err error
	if !s.file.CheckFileExists(oracleObj.File.FileDetails.Path) && oracleObj.File.FileDetails.GitDetails.Action != domain.DeleteAction {
		s.writeError(oracleObj, domain.FileNotExists)
	}

	parts := strings.Split(oracleObj.File.FileDetails.Path, string(os.PathSeparator))
	if len(parts) >= 4 {
		oracleObj.EpicModuleName = parts[len(parts)-4]
		oracleObj.ModuleName = parts[len(parts)-3]
		oracleObj.ObjectName = strings.ToLower(oracleObj.File.FileDetails.Name[:len(oracleObj.File.FileDetails.Name)-len(filepath.Ext(oracleObj.File.FileDetails.Name))])

		oracleObj.ObjectType, err = s.getObjectTypeFromDir(parts[len(parts)-2], oracleObj.ObjectName)
		if err != nil {
			oracleObj.Errors = append(oracleObj.Errors, err.Error())
		}
	} else {
		s.writeError(oracleObj, domain.UnknownObjectType)
	}
}

func (s *service) addSchema(oracleObj *domain.OracleObject) {
	serverSchemaAliasLine, err := s.file.SearchStrInFile("schema", oracleObj.File.FileDetails.Path)
	if err != nil {
		s.writeError(oracleObj, domain.SchemaNotFound)
	}

	for key := range s.schema.ParseSchemaAliasLine(serverSchemaAliasLine) {
		oracleObj.ServerSchemaList = append(oracleObj.ServerSchemaList, key)
	}

	if len(oracleObj.ServerSchemaList) == 0 {
		s.writeError(oracleObj, domain.SchemaNotFound)
	}
}

func (s *service) CreateOracleObjects(rootDir, installDir string, files []domain.File) []domain.OracleObject {
	result := make([]domain.OracleObject, 0, len(files))
	for _, f := range files {
		if s.fileSuitable(rootDir, installDir, f.Path) {
			oracleFile := domain.OracleFile{
				OracleDataType: domain.Data,
				FileDetails:    f,
			}
			obj := domain.OracleObject{File: oracleFile}
			s.resolveAdditionalPathInfo(&obj)
			if f.GitDetails.Action != domain.DeleteAction {
				s.addSchema(&obj)
			}
			result = append(result, obj)
		}
	}

	return result
}

package extractor

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/filewalker"
	"os"
	"path/filepath"
	"strings"
)

type Service interface {
	ExtractOracleObjects(files []domain.File) []domain.OracleObject
	WalkAndExtractOracleObjects(rootDir string) ([]domain.OracleObject, error)
}

type service struct {
	fileWalker filewalker.Service
}

var _ Service = (*service)(nil)

func NewService(fileWalker filewalker.Service) *service {
	return &service{fileWalker: fileWalker}
}

func (s *service) getObjectTypeFromDir(objectTypeDir, objectName string) (int, error) {
	matchMap := map[string]int{
		"tablespaces": domain.OracleTablespaceType,
		"directories": domain.OracleDirectoryType,
		"dblinks":     domain.OracleDbLinkType,
		"users":       domain.OracleUserType,
		"synonyms":    domain.OracleSynonymType,
		"contexts":    domain.OracleContextType,
		"sequences":   domain.OracleSequencesType,
		"types":       domain.OracleTypeType,
		"tables":      domain.OracleTableType,
		"mlogs":       domain.OracleMLogType,
		"mviews":      domain.OracleMViewType,
		"packages":    domain.OraclePackageType,
		"views":       domain.OracleViewType,
		"triggers":    domain.OracleTriggerType,
		"vtbs_tasks":  domain.OracleVTaskType,
		"rows":        domain.OracleRowType,
		"roles":       domain.OracleRoleType,
		"functions":   domain.OracleFunctionType,
		"vtbs_clogs":  domain.OracleVClogType,
		"tables.fk":   domain.OracleTableFKType,
	}
	if val, ok := matchMap[objectTypeDir]; ok {
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

func (s *service) addServerSchema(headersMap map[domain.ServerSchema]struct{}, headerStr string) {
	if headerStr == "core" {
		headersMap[domain.ServerSchema{
			Server: "core",
			Schema: "vtbs",
		}] = struct{}{}
	} else if headerStr == "charger" || headerStr == "hpffm" {
		headersMap[domain.ServerSchema{
			Server: "hpffm",
			Schema: "vtbs",
		}] = struct{}{}
	} else if headerStr == "vtbs_bi" {
		headersMap[domain.ServerSchema{
			Server: "hpffm",
			Schema: "vtbs_bi",
		}] = struct{}{}
	} else if headerStr == "vtbs_x_alaris" || headerStr == "xalaris" {
		headersMap[domain.ServerSchema{
			Server: "hpffm",
			Schema: "vtbs_x_alaris",
		}] = struct{}{}
	} else if headerStr == "adesk" || headerStr == "vtbs_adesk" || headerStr == "reporter" {
		headersMap[domain.ServerSchema{
			Server: "hpffm",
			Schema: "vtbs_adesk",
		}] = struct{}{}
	}
}

func (s *service) parseSchema(schemaStr string) map[domain.ServerSchema]struct{} {
	schemaStr = strings.ToLower(schemaStr)
	schemaStr = strings.Replace(schemaStr, "schema", "", -1)
	schemaStr = strings.Replace(schemaStr, ":", "", -1)
	schemaStr = strings.Replace(schemaStr, " ", "", -1)
	schemaStr = strings.Replace(schemaStr, "--", "", -1)
	schemaStr = strings.Replace(schemaStr, "/", ",", -1)
	schemaStr = strings.Replace(schemaStr, "\\", ",", -1)
	if len(schemaStr) > 0 {
		parts := strings.Split(schemaStr, ",")
		result := make(map[domain.ServerSchema]struct{})
		for _, val := range parts {
			s.addServerSchema(result, val)
		}
		return result
	}
	return nil
}

func (s *service) ExtractOracleObjects(files []domain.File) []domain.OracleObject {
	var err error
	result := make([]domain.OracleObject, 0, len(files))
	for _, file := range files {
		obj := domain.OracleObject{File: file}
		if !s.fileWalker.CheckFileExists(file.Path) && file.GitAction != domain.DeleteAction {
			s.writeError(&obj, domain.FileNotExists)
		}

		parts := strings.Split(file.Path, string(os.PathSeparator))
		if len(parts) >= 4 {
			obj.EpicModuleName = parts[len(parts)-4]
			obj.ModuleName = parts[len(parts)-3]
			obj.ObjectName = strings.ToLower(file.Name[:len(file.Name)-len(filepath.Ext(file.Name))])
			obj.ObjectType, err = s.getObjectTypeFromDir(parts[len(parts)-2], obj.ObjectName)
			if err != nil {
				obj.Errors = append(obj.Errors, err.Error())
			}
		} else {
			s.writeError(&obj, domain.UnknownObjectType)
		}

		serverSchema, err := s.fileWalker.SearchStrInFile("schema", file.Path)
		if err != nil {
			s.writeError(&obj, domain.SchemaNotFound)
		}

		for key := range s.parseSchema(serverSchema) {
			obj.ServerSchemaList = append(obj.ServerSchemaList, key)
		}

		if len(obj.ServerSchemaList) == 0 {
			s.writeError(&obj, domain.SchemaNotFound)
		}
		result = append(result, obj)
	}

	return result
}

func (s *service) WalkAndExtractOracleObjects(rootDir string) ([]domain.OracleObject, error) {
	files, err := s.fileWalker.Walk(rootDir, domain.OracleFileType)
	if err != nil {
		return nil, err
	}
	return s.ExtractOracleObjects(files), nil
}

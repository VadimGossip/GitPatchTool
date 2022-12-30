package extractor

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/filewalker"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type Service interface {
	ExtractOracleObjects(files []domain.File) []domain.OracleObject
}

type service struct {
	fileWalker filewalker.Service
}

var _ Service = (*service)(nil)

func NewService(fileWalker filewalker.Service) *service {
	return &service{fileWalker: fileWalker}
}

func (s *service) getObjectTypeFromDir(objectTypeDir string) (int, error) {
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
	}
	if val, ok := matchMap[objectTypeDir]; ok {
		return val, nil
	}

	return 0, domain.UnknownObjectType
}

func (s *service) writeError(obj *domain.OracleObject, err error) {
	if err != nil {
		obj.Errors = append(obj.Errors, err.Error())
	}
}

func (s *service) addHeader(headersMap map[objectHeader]struct{}, headerStr string) {
	if headerStr == "core" {
		headersMap[objectHeader{
			server: "core",
			schema: "vtbs",
		}] = struct{}{}
	} else if headerStr == "charger" || headerStr == "hpffm" {
		headersMap[objectHeader{
			server: "hpffm",
			schema: "vtbs",
		}] = struct{}{}
	} else if headerStr == "vtbs_bi" {
		headersMap[objectHeader{
			server: "hpffm",
			schema: "vtbs_bi",
		}] = struct{}{}
	} else if headerStr == "vtbs_x_alaris" || headerStr == "xalaris" {
		headersMap[objectHeader{
			server: "hpffm",
			schema: "vtbs_x_alaris",
		}] = struct{}{}
	} else if headerStr == "adesk" || headerStr == "vtbs_adesk" || headerStr == "reporter" {
		headersMap[objectHeader{
			server: "hpffm",
			schema: "vtbs_adesk",
		}] = struct{}{}
	}
}

func (s *service) parseSchema(schemaStr string) map[objectHeader]struct{} {
	schemaStr = strings.ToLower(schemaStr)
	schemaStr = strings.Replace(schemaStr, "schema", "", -1)
	schemaStr = strings.Replace(schemaStr, ":", "", -1)
	schemaStr = strings.Replace(schemaStr, " ", "", -1)
	schemaStr = strings.Replace(schemaStr, "--", "", -1)
	schemaStr = strings.Replace(schemaStr, "/", ",", -1)
	schemaStr = strings.Replace(schemaStr, "\\", ",", -1)
	if len(schemaStr) > 0 {
		parts := strings.Split(schemaStr, ",")
		result := make(map[objectHeader]struct{})
		for _, val := range parts {
			s.addHeader(result, val)
		}
		return result
	}
	return nil
}

func (s *service) ExtractOracleObjects(files []domain.File) []domain.OracleObject {
	result := make([]domain.OracleObject, 0, len(files))
	for _, file := range files {
		obj := domain.OracleObject{File: file}
		if !s.fileWalker.CheckFileExists(file.Path) && file.GitAction != domain.DeleteAction {
			s.writeError(&obj, domain.FileNotExists)
			result = append(result, obj)
			continue
		}

		schema, err := s.fileWalker.SearchStrInFile("schema", file.Path)
		if err != nil {
			s.writeError(&obj, domain.FileNotExists)
			result = append(result, obj)
			continue
		}
		logrus.Info(s.parseSchema(schema))

		parts := strings.Split(file.ShortPath, string(os.PathSeparator))
		if len(parts) < 4 {
			s.writeError(&obj, domain.UnknownObjectType)
			result = append(result, obj)
			continue
		}
		obj.EpicModuleName = parts[len(parts)-4]
		obj.ModuleName = parts[len(parts)-3]
		obj.ObjectName = file.Name[:len(file.Name)-len(filepath.Ext(file.Name))]
		obj.ObjectType, err = s.getObjectTypeFromDir(parts[len(parts)-2])
		if err != nil {
			obj.Errors = append(obj.Errors, err.Error())
		}
		result = append(result, obj)
	}

	return result
}

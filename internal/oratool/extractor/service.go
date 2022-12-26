package extractor

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/filewalker"
	"os"
	"path/filepath"
	"strings"
)

type Service interface {
	ExtractOracleObjects(rootDir string, files []domain.File) []domain.OracleObject
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

func (s *service) ExtractOracleObjects(rootDir string, files []domain.File) []domain.OracleObject {
	var err error
	result := make([]domain.OracleObject, 0, len(files))
	for _, file := range files {
		obj := domain.OracleObject{File: file}
		if !s.fileWalker.CheckFileExists(rootDir+file.Path) && file.GitAction != domain.DeleteAction {
			s.writeError(&obj, domain.FileNotExists)
			result = append(result, obj)
			continue
		}

		parts := strings.Split(file.Path, string(os.PathSeparator))
		if len(parts) > 0 && parts[0] == "install" {
			continue
		}

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

package extractor

import (
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"os"
	"strings"
)

type Service interface {
	ExtractOracleObjects(files []domain.File) []domain.OracleObject
}

type service struct {
}

var _ Service = (*service)(nil)

func NewService() *service {
	return &service{}
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

func (s *service) ExtractOracleObjects(files []domain.File) []domain.OracleObject {
	var err error
	result := make([]domain.OracleObject, 0, len(files))
	for _, file := range files {
		parts := strings.Split(file.Path, string(os.PathSeparator))
		if len(parts) >= 4 {
			obj := domain.OracleObject{
				EpicModuleName: parts[len(parts)-4],
				ModuleName:     parts[len(parts)-3],
				ObjectType:     0,
				Schema:         "",
				Server:         "",
				File:           file,
			}
			obj.ObjectType, err = s.getObjectTypeFromDir(parts[len(parts)-2])
			if err != nil {
				obj.Errors = append(obj.Errors, err.Error())
			}
			result = append(result, obj)
		}
	}

	return result
}

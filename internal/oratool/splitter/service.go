package splitter

import (
	"bufio"
	"errors"
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/filewalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"
)

type Service interface {
	SplitTableFiles() error
}

type service struct {
	cfg        *domain.Config
	fileWalker filewalker.Service
	extractor  extractor.Service
}

var _ Service = (*service)(nil)

func NewService(cfg *domain.Config, fileWalker filewalker.Service, extractor extractor.Service) *service {
	return &service{cfg: cfg, fileWalker: fileWalker, extractor: extractor}
}

type file struct {
	name string
	text []string
}

func (s *service) markFileLines(fileLines []string) map[int]string {
	stopElement := ";"

	foundLines := make([]string, 0)
	markedMap := make(map[int]string)

	for idx, fileLine := range fileLines {
		tmpFileLine := strings.TrimSpace(strings.ToLower(fileLine))
		tmpFileLine = regexp.MustCompile(`\s+`).ReplaceAllString(tmpFileLine, " ")

		if strings.HasPrefix(tmpFileLine, "alter") {
			foundLines = []string{tmpFileLine}
		} else {
			foundLines = append(foundLines, tmpFileLine)
		}

		resultStr := ""
		for _, str := range foundLines {
			if resultStr == "" {
				resultStr = strings.TrimSpace(str)
			} else {
				resultStr = resultStr + " " + strings.TrimSpace(str)
			}
		}

		parts := strings.Split(resultStr, " ")
		if len(parts) == 12 &&
			parts[0] == "alter" &&
			parts[1] == "table" &&
			parts[3] == "add" &&
			parts[4] == "constraint" &&
			parts[6] == "foreign" &&
			parts[7] == "key" &&
			parts[9] == "references" && strings.HasSuffix(parts[len(parts)-1], stopElement) {
			for i := idx - (len(foundLines) - 1); i <= idx; i++ {
				markedMap[i] = parts[5]
			}
			foundLines = nil
		}
	}
	return markedMap
}

func (s *service) createDir(dirPath string) error {
	if _, err := os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
		if err = os.Mkdir(dirPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) createFile(filePath string, fileLines []string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, fileLine := range fileLines {
		_, err := f.WriteString(fileLine + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) splitTableFile(oraFile domain.OracleObject) error {
	f, err := os.Open(oraFile.File.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	fileLines := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}

	markedLines := s.markFileLines(fileLines)

	filesToCreate := make(map[string][]string)
	for idx, fileLine := range fileLines {
		if val, ok := markedLines[idx]; ok {
			if _, ok := filesToCreate[val+".sql"]; !ok {
				schema, err := s.fileWalker.SearchStrInFile("schema", oraFile.File.Path)
				if err == nil {
					filesToCreate[val+".sql"] = append(filesToCreate[val], schema)
				}
			}
			filesToCreate[val+".sql"] = append(filesToCreate[val+".sql"], fileLine)

		} else {
			filesToCreate[oraFile.File.Name] = append(filesToCreate[oraFile.File.Name], fileLine)

		}
	}
	path := oraFile.File.Path

	if len(filesToCreate) > 1 {
		if err := s.createDir(path[:len(oraFile.File.Path)-len("tables"+string(os.PathSeparator)+oraFile.File.Name)] + "tables.fk"); err != nil {
			return err
		}
		for key, val := range filesToCreate {
			if key != oraFile.File.Name {
				path = path[:len(oraFile.File.Path)-len("tables"+string(os.PathSeparator)+oraFile.File.Name)] + "tables.fk" + string(os.PathSeparator) + oraFile.ObjectName + "." + key

			}
			if err := s.createFile(path, val); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *service) SplitTableFiles() error {
	oraObjects, err := s.extractor.WalkAndExtractOracleObjects(s.cfg.Path.RootDir)
	if err != nil {
		return err
	}

	filteredObj := make([]domain.OracleObject, 0)
	for _, val := range oraObjects {
		if val.ObjectType == domain.OracleTableType && val.EpicModuleName == "stlm" && val.ModuleName == "Settlements" {
			filteredObj = append(filteredObj, val)
		}
	}

	logrus.Info(len(filteredObj))

	if err := s.splitTableFile(filteredObj[1]); err != nil {
		return err
	}
	return nil
}

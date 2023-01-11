package splitter

// splitter need to refactor service to be more flexible

import (
	"bufio"
	"errors"
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/file"
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
	fileWalker file.Service
	extractor  extractor.Service
}

var _ Service = (*service)(nil)

func NewService(cfg *domain.Config, fileWalker file.Service, extractor extractor.Service) *service {
	return &service{cfg: cfg, fileWalker: fileWalker, extractor: extractor}
}

func (s *service) markFileLines(fileLines []string) map[int]string {
	foundLines := make([]string, 0)
	markedMap := make(map[int]string)
	sameName := make(map[string]struct{})

	for idx, fileLine := range fileLines {
		tmpFileLine := strings.TrimSpace(strings.ToLower(fileLine))
		tmpFileLine = regexp.MustCompile(`\s+`).ReplaceAllString(tmpFileLine, " ")

		if strings.HasPrefix(tmpFileLine, "alter") || strings.HasPrefix(tmpFileLine, "create") {
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

		//fmt.Printf("index %d tmpFileLine %s\n", idx, tmpFileLine)

		parts := strings.Split(resultStr, " ")
		if len(parts) >= 11 &&
			parts[0] == "alter" &&
			parts[1] == "table" &&
			parts[3] == "add" &&
			(parts[4] == "constraint" || parts[5] == "constraint") &&
			(parts[6] == "foreign" || parts[7] == "foreign") &&
			(parts[7] == "key" || parts[8] == "key") &&
			(parts[9] == "references" || parts[10] == "references") &&
			(strings.HasSuffix(parts[len(parts)-1], ";") || strings.HasSuffix(parts[len(parts)-1], "/")) {

			//fmt.Println(resultStr)

			var name string
			for k := range parts {
				if parts[k] == "constraint" {
					name = strings.Replace(parts[k+1], `"`, "", -1)
					break
				}
			}
			for i := idx - (len(foundLines) - 1); i <= idx; i++ {
				markedMap[i] = name
			}
			foundLines = nil
			if len(parts) > 16 {
				logrus.Infof("Split manual %s", name)
			}
			if _, ok := sameName[name]; !ok {
				sameName[name] = struct{}{}
			} else {
				logrus.Infof("Same name found name %s", name)
			}
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
			filesToCreate[val+".sql"] = append(filesToCreate[val+".sql"], strings.ToLower(fileLine))

		} else {
			filesToCreate[oraFile.File.Name] = append(filesToCreate[oraFile.File.Name], fileLine)

		}
	}

	if len(filesToCreate) > 1 {
		if err := s.createDir(oraFile.File.Path[:len(oraFile.File.Path)-len("tables"+string(os.PathSeparator)+oraFile.File.Name)] + "tables.fk"); err != nil {
			return err
		}
		for key, val := range filesToCreate {
			path := oraFile.File.Path
			if key != oraFile.File.Name {
				path = oraFile.File.Path[:len(oraFile.File.Path)-len("tables"+string(os.PathSeparator)+oraFile.File.Name)] + "tables.fk" + string(os.PathSeparator) + oraFile.ObjectName + "." + key
			}

			//fmt.Printf("key %s %d\n", path, len(val))

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
		if val.ObjectType == domain.OracleTableType {
			filteredObj = append(filteredObj, val)
		}
	}

	for _, item := range filteredObj {
		if err := s.splitTableFile(item); err != nil {
			return err
		}
	}

	return nil
}

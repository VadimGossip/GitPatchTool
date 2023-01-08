package splitter

import (
	"bufio"
	"fmt"
	"github.com/VadimGossip/gitPatchTool/internal/domain"
	"github.com/VadimGossip/gitPatchTool/internal/filewalker"
	"github.com/VadimGossip/gitPatchTool/internal/oratool/extractor"
	"github.com/sirupsen/logrus"
	"os"
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

func (s *service) markFileLines(fileLines []string) (map[string][]int, error) {
	stopElement := ";"
	searchTmpl := []string{"create", "any", "index", "any", "on", "any", "any", "any"}
	foundLines := make([]string, 0)
	markedMap := make(map[string][]int)

	for idx, fileLine := range fileLines {
		tmpFileLine := strings.TrimSpace(strings.ToLower(fileLine))

		if strings.HasPrefix(tmpFileLine, searchTmpl[0]) {
			foundLines = []string{tmpFileLine}
		} else {
			foundLines = append(foundLines, tmpFileLine)
		}

		resultStr := ""
		for _, str := range foundLines {
			if resultStr == "" {
				resultStr = strings.TrimSpace(str)
			}
			resultStr += " " + strings.TrimSpace(str)
		}

		parts := strings.Split(resultStr, " ")
		if len(parts) > 4 && (parts[2] == "index" || parts[1] == "index") && strings.HasSuffix(parts[len(parts)-1], stopElement) {
			for i := idx; i < idx+len(foundLines); i++ {
				markedMap[parts[3]] = append(markedMap[parts[3]], i)
			}
			fmt.Println(resultStr)
			foundLines = nil
		}
	}
	return markedMap, nil
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

	logrus.Infof("len filteredObj %d", len(filteredObj))

	logrus.Info(filteredObj[1].File.Path)

	f, err := os.Open(filteredObj[1].File.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	fileLines := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		//fileLine := strings.TrimSpace(strings.ToLower(scanner.Text()))
		//fileLine = regexp.MustCompile(`\s+`).ReplaceAllString(fileLine, " ")
		fileLines = append(fileLines, scanner.Text())
	}

	res, err := s.markFileLines(fileLines)
	fmt.Println(res)
	return nil
}

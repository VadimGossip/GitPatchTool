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

func (s *service) markLinesByTmpl(startPos int, currentLine string, foundLines *[]string, markedMap map[string][]int) error {
	stopElement := ";"
	searchTmpl := []string{"create", "any", "index", "any", "on", "any", "any", "any"}

	if strings.HasPrefix(currentLine, searchTmpl[0]) {
		foundLines = &[]string{currentLine}
	} else {
		*foundLines = append(*foundLines, currentLine)
	}

	resultStr := ""
	for _, str := range *foundLines {
		if resultStr == "" {
			resultStr = strings.TrimSpace(str)
		}
		resultStr += " " + strings.TrimSpace(str)
	}

	parts := strings.Split(resultStr, " ")
	if (parts[2] == "index" || parts[1] == "index") && strings.HasSuffix(parts[len(parts)-1], stopElement) {

		for i := startPos; i <= len(*foundLines); i++ {
			markedMap[parts[3]] = append(markedMap[parts[3]], i)
		}
		fmt.Println(resultStr)
		foundLines = &[]string{}

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

	//	searchTmpl := []string{"create", "any", "index", "any", "on", "any", "any", "any", ";"}

	foundLines := make([]string, 0)
	markedMap := make(map[string][]int)
	for idx, fileLine := range fileLines {
		tmpFileLine := strings.TrimSpace(strings.ToLower(fileLine))
		s.markLinesByTmpl(idx, tmpFileLine, &foundLines, markedMap)
		for _, x := range foundLines {
			fmt.Println(x)
		}

	}

	//for idx, fileLine := range fileLines {
	//	//tmpFileLine := regexp.MustCompile(`\s+`).ReplaceAllString(fileLine, " ")
	//	//fmt.Println(tmpFileLine)
	//	tmpFileLine := strings.TrimSpace(strings.ToLower(fileLine))
	//	s.markLinesByTmpl(idx, tmpFileLine, &foundLines, markedMap)
	//	////if len(foundLines) > 0 {
	//	////   for _, str := range foundLines{
	//	////	   if tmpStr == ""{
	//	////		   tmpStr = str
	//	////	   }
	//	////	   tmpStr += " "+str
	//	////   }
	//	////}
	//	////
	//	////parts := strings.Split(tmpFileLine, string(os.PathSeparator))
	//	//
	//	//if len(foundLines) > 0 && !strings.HasPrefix(tmpFileLine, searchTmpl[0]) ||
	//	//	len(foundLines) == 0 && strings.HasPrefix(tmpFileLine, searchTmpl[0]) {
	//	//
	//	//	foundLines = append(foundLines, tmpFileLine)
	//	//	resultStr = ""
	//	//	for _, str := range foundLines {
	//	//		if resultStr == "" {
	//	//			resultStr = strings.TrimSpace(str)
	//	//		}
	//	//		resultStr += " " + strings.TrimSpace(str)
	//	//	}
	//	//
	//	//	parts := strings.Split(resultStr, " ")
	//	//	for idx, val := range parts {
	//	//		fmt.Printf("idx %d val %s\n", idx, val)
	//	//	}
	//	//	fmt.Printf("last %s\n ", parts[len(parts)-1])
	//	//	if (parts[2] == "index" || parts[1] == "index") && strings.HasSuffix(parts[len(parts)-1], ";") {
	//	//		foundLines = nil
	//	//		fmt.Println("Nice")
	//	//	}
	//	//
	//	//	//fmt.Println(foundLines)
	//	//	//|| parts[2] == "index") && strings.HasSuffix(parts[len(parts)-1], ";")
	//	//	//for _, v := range parts {
	//	//	//	fmt.Println(v)
	//	//	//}
	//	//	//
	//	//	//fmt.Println()
	//	//	//fmt.Println(parts)
	//	//	//fmt.Printf("last of parts %s", parts[len(parts)-1])
	//	//	//fmt.Println()
	//	//
	//	//	//if parts[2] == "index" || parts[1] == "index" {
	//	//	//	//fmt.Printf("last of parts %s", parts[len(parts)-1])
	//	//	//	fmt.Println(resultStr)
	//	//	//} else {
	//	//	//
	//	//	//}
	//	//
	//	//} else {
	//	//	foundLines = nil
	//	//}
	//
	//	//if strings.HasPrefix(tmpFileLine, searchTmpl[0]) {
	//	//	if len(foundLines) == 0 {
	//	//		foundLines = append(foundLines, fileLine)
	//	//		if tmpFileLine[len[tmpFileLine]-1] == searchTmpl[0]
	//	//
	//	//	} else {
	//	//		foundLines = nil
	//	//	}
	//	//}
	//
	//	//if len(foundLines) == 3
	//	//
	//	//if len(searchLines) == 0 && strings.HasPrefix(fileLines[idx], searchTmpl[0]) {
	//	//	searchLines = append()
	//	//} else if len(searchLines) > 0 {
	//	//	if strings.HasSuffix(fileLines[idx], searchTmpl[len(searchTmpl) - 1]) {
	//	//
	//	//	}
	//	//}
	//
	//	//for sIdx := range searchTmpl {
	//	//	if idx + sIdx + 1 > len(fileLines) {
	//	//		break
	//	//	}
	//	//	if strings.HasPrefix(fileLines[idx + sIdx], searchTmpl[sIdx]) {
	//	//
	//	//	}
	//	//}
	//	//if strings.HasPrefix(fileLine, searchTmpl[0]) {
	//	//	newLine
	//	//}
	//}

	return nil
}

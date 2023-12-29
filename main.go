package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) <= 2 {
		log.Fatalln("Too few arguments.")
		return
	}
	dir := os.Args[1]
	patternList := os.Args[2:]
	importsMap := make(map[string]int)

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("訪問中のエラー:", err)
			return err
		}
		if strings.HasSuffix(path, ".go") {
			var lines []string
			lines, err = parseFile(path, patternList, importsMap)
			if err == nil && len(lines) > 0 {
				fmt.Println("ファイル: ", path)
				for _, line := range lines {
					fmt.Println(line)
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("エラーが発生しました:", err)
		return
	}
	pairs := make([]struct {
		Key   string
		Value int
	}, 0, len(importsMap))

	for k, v := range importsMap {
		pairs = append(pairs, struct {
			Key   string
			Value int
		}{k, v})
	}

	// 値に基づいてスライスをソート
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	for _, pair := range pairs {
		fmt.Printf("import: %s, 回数: %d\n", pair.Key, pair.Value)
	}
}

func parseFile(filePath string, patternList []string, importsMap map[string]int) ([]string, error) {
	cnt := 0
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	imports := []string{}
	if err != nil {
		fmt.Println("ファイル解析中にエラーが発生しました:", err)
		return imports, err
	}

	for _, i := range node.Imports {
		importPath := strings.Trim(i.Path.Value, `"`)
		if allMatch(importPath, patternList) {
			importsMap[importPath]++
			imports = append(imports, "  import: "+importPath)
			cnt++
		}

	}
	return imports, nil
}

func oneOfMatch(w string, patternList []string) bool {
	if len(patternList) <= 0 {
		return false
	}

	return strings.Contains(w, patternList[0]) || oneOfMatch(w, patternList[1:])
}

func allMatch(w string, patternList []string) bool {
	if len(patternList) <= 0 {
		return true
	}

	return strings.Contains(w, patternList[0]) && allMatch(w, patternList[1:])
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

const folderPath = ""

func main() {
	var todayResults []Properties
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			props, err := extractFrontmatterProperties(path)
			if err != nil {
				fmt.Printf("Error parsing file %s: %v\n", path, err)
			}
			if props == nil {
				return nil
			} else {
				todayResults = append(todayResults, *props)
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println("Error walking the directory:", err)
	}

	writeCSVUpdateIfNeeded(time.Now().Format("2006-01-02"), todayResults)
}

type Properties struct {
	ReportKey      string
	TomatoesPerDay int
	Status         string
	Priority       int
	DoneTotal      int
	Done           map[string]int
}

func extractFrontmatterProperties(filePath string) (*Properties, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	content := string(data)
	lines := strings.SplitN(content, "---", 3) // Split at the frontmatter delimiters

	if len(lines) < 3 {
		// Not a project file
		return nil, nil
	}

	frontMatter := lines[1]
	// markdownContent := lines[2]

	var fm map[string]string
	if err := yaml.Unmarshal([]byte(frontMatter), &fm); err != nil {
		// Not a project file
		fmt.Printf("not Formatted: %v\n", filePath)
		return nil, nil
	}

	shouldExist(filePath, fm, "reportKey")

	return mapToProperties(fm), nil
}

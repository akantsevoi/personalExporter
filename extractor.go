package main

import (
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// I report everything in tomatoes, and the time per tomato may wary
type Properties struct {
	ReportKey      string
	TomatoesPerDay int
	Status         string
	Priority       int
	DoneTotal      int
	Done           map[string]int
}

// checks if file has FrontMatter properties
// if no - returns nil
// if yes - checks for mandatory fields and returns properties
func extractFrontmatterProperties(filePath string) (*Properties, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	content := string(data)
	// Check if has FrontMatter delimiter
	if !strings.HasPrefix(content, "---") {
		return nil, nil
	}
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

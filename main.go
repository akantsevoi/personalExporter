package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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
			props, err := parseMarkdownFile(path)
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

	sort.Slice(todayResults, func(i, j int) bool {
		return todayResults[i].ReportKey > todayResults[j].ReportKey
	})

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

func parseMarkdownFile(filePath string) (*Properties, error) {
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

func mapToProperties(input map[string]string) *Properties {
	suffix := "_sub_done_today"

	pr := Properties{
		ReportKey:      input["reportKey"],
		Priority:       toInt(input["priority"]),
		Status:         input["status"],
		TomatoesPerDay: toInt(input["tomatosPerDay"]),
		Done:           map[string]int{},
	}

	if v, ok := input["doneToday"]; ok {
		pr.DoneTotal = toInt(v)
	} else {
		for k, v := range input {
			name, hasSuffix := strings.CutSuffix(k, suffix)
			if !hasSuffix {
				continue
			}

			// log.Printf("%v %T %v\n", k, v, v)

			increment := toInt(v)
			pr.DoneTotal += increment
			pr.Done[name] = increment
		}
	}

	return &pr
}

func toInt(v string) int {
	i, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Errorf("parseInt: %v: %w", v, err))
	}
	return i
}

func shouldExist(path string, m map[string]string, key string) {
	if _, ok := m[key]; !ok {
		fmt.Printf("no %v in %v\n%v", key, path, m)
		os.Exit(1)
	}
}

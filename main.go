package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	folderPath     = ""
	hoursPerTomato = float64(0.0)
)

func main() {

	sourcePath := flag.String("source", "", "path to the root folder of projects")
	hoursPT := flag.Float64("hpert", 1.0, "hours per reported tomato")
	writeFolder := flag.String("result", "", "where to save the result\nIf empty - it will be stored in the same folder as source")

	flag.Parse()

	if len(*writeFolder) == 0 {
		writeFolder = sourcePath
	}

	fmt.Println(*sourcePath)
	fmt.Println(*writeFolder)
	fmt.Println(*hoursPT)

	folderPath = *sourcePath
	hoursPerTomato = *hoursPT

	// work
	paths, err := markdownFilePaths(folderPath)
	stopIfErrf("%w", err)

	todayResults, err := extractPropsFromFiles(paths)
	stopIfErrf("%w", err)

	writeCSVUpdateIfNeeded(time.Now().Format("2006-01-02"), *writeFolder, todayResults)
}

func extractPropsFromFiles(filePaths []string) ([]Properties, error) {
	var results []Properties
	for _, p := range filePaths {
		props, err := extractFrontmatterProperties(p)
		if err != nil {
			return nil, fmt.Errorf("extractPropsFromFiles: error parsing file %s: %w", p, err)
		}
		if props == nil {
			continue
		}

		results = append(results, *props)
	}
	return results, nil
}

func markdownFilePaths(projectsFolder string) ([]string, error) {
	var fileNames []string
	err := filepath.Walk(projectsFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			fileNames = append(fileNames, path)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("markdownFilePaths: error walking the directory: %w", err)
	}

	return fileNames, nil
}

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

var (
	hoursPerTomato = float64(0.0)
)

type Command string

const (
	commandReport Command = "report"
	commandReset  Command = "reset"
	commandStat   Command = "stat"
)

var allCommands = []Command{
	commandReport,
	commandReset,
	commandStat,
}

func main() {

	commandStr := flag.String("command", "", fmt.Sprintf("what to do {%v}", allCommands))
	sourcePath := flag.String("source", "", "path to the root folder of projects")
	hoursPT := flag.Float64("hpert", 1.0, "hours per reported tomato")
	writeFolder := flag.String("result", "", "where to save the result\nIf empty - it will be stored in the same folder as source")

	// stat
	lastDays := flag.Int("lDays", 7, "last days statistics")
	statsFile := flag.String("sFile", "", "path to statistics file")

	flag.Parse()

	if len(*commandStr) == 0 {
		fmt.Println("command is mandatoty")
		os.Exit(1)
	}

	command := Command(*commandStr)
	if !slices.Contains(allCommands, command) {
		fmt.Println("unsupported command: ", command)
		os.Exit(1)
	}

	switch command {
	case commandReport:
		folderPath := *sourcePath
		paths, err := markdownFilePaths(folderPath)
		stopIfErrf("%w", err)

		if len(*writeFolder) == 0 {
			writeFolder = sourcePath
		}
		hoursPerTomato = *hoursPT

		fmt.Println("report start...")
		fmt.Println("source:\t\t", folderPath)
		fmt.Println("result:\t\t", *writeFolder)
		fmt.Println("hours/tomato:\t", hoursPerTomato)

		todayResults, err := extractPropsFromFiles(paths)
		stopIfErrf("report: %w", err)

		writeCSVUpdateIfNeeded(time.Now().Format("2006-01-02"), *writeFolder, todayResults)
	case commandReset:
		folderPath := *sourcePath
		paths, err := markdownFilePaths(folderPath)
		stopIfErrf("%w", err)

		fmt.Println("reset start...")
		fmt.Println("project files:\t", folderPath)
		err = resetProgress(paths)
		stopIfErrf("reset: %w", err)
	case commandStat:
		if len(*statsFile) == 0 {
			fmt.Println("sFile is empty")
			os.Exit(1)
		}
		err := stats(*statsFile, *lastDays)
		stopIfErrf("stat: %w", err)
	}

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

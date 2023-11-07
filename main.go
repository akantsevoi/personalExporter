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
	folderPath      = ""
	writeFolderPath = ""
	hoursPerTomato  = float64(0.0)
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
	writeFolderPath = *writeFolder
	hoursPerTomato = *hoursPT

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

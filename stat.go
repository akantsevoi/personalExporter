package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strings"
)

// ignores last day as not relevant(incomplete)
func stats(resultFilePath string, lastDays int) error {
	originalFile, err := os.OpenFile(resultFilePath, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0644)
	stopIfErrf("open originalfile err: %w", err)

	reader := csv.NewReader(originalFile)
	fileContent, err := reader.ReadAll()
	stopIfErrf("read csv lines: %#v", err)

	var projectIndexes []int
	for i, header := range fileContent[0] {
		if strings.Contains(header, ":") || i == 0 {
			continue
		}
		projectIndexes = append(projectIndexes, i)
	}

	values := make([]float64, 0, lastDays)
	var sum float64
	for i := len(fileContent) - 1 - lastDays - 1; i < len(fileContent)-1; i++ {
		var totalDay float64
		for _, pjIndex := range projectIndexes {
			totalDay += toFloat64(fileContent[i][pjIndex])
		}
		fmt.Println(fileContent[i][0], " - ", totalDay)
		sum += totalDay
		values = append(values, totalDay)
	}

	slices.Sort(values)
	var median float64
	if len(values)%2 == 0 {
		a, b := values[len(values)/2-1], values[len(values)/2]
		median = (a + b) / 2
	} else {
		median = values[len(values)/2]
	}

	fmt.Println("mean:\t", sum/float64(len(values)))
	fmt.Println("median:\t", median)

	return nil
}

func projectStats(resultFilePath string, projectName string) error {
	originalFile, err := os.OpenFile(resultFilePath, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0644)
	stopIfErrf("open originalfile err: %w", err)

	reader := csv.NewReader(originalFile)
	fileContent, err := reader.ReadAll()
	stopIfErrf("read csv lines: %#v", err)

	var indexes columnIndexes
	if ind := projectIndexes(fileContent[0], projectName); ind != nil {
		indexes = *ind
	} else {
		return fmt.Errorf("no such project or format error")
	}

	total := 0.0
	perSubTotal := []pair[string, float64]{}
	for _, sp := range indexes.subProjs {
		perSubTotal = append(perSubTotal, pair[string, float64]{
			sp.First, 0.0,
		})
	}

	for ri := 1; ri < len(fileContent); ri++ {
		total += toFloat64(fileContent[ri][indexes.proj])

		var resIndex int
		for _, sp := range indexes.subProjs {
			perSubTotal[resIndex].Second += toFloat64(fileContent[ri][sp.Second])
			resIndex++
		}
	}

	fmt.Println("proj:\t", projectName)
	fmt.Println("total:\t", total)
	for _, sT := range perSubTotal {
		fmt.Println("\t", sT.First, "\t", sT.Second)
	}

	return nil
}

type columnIndexes struct {
	proj     int
	subProjs []pair[string, int]
}

// returns nil if there is no such project
func projectIndexes(headerRow []string, projectName string) *columnIndexes {
	if len(headerRow) == 0 || len(projectName) == 0 {
		return nil
	}

	mainProjectIndex := -1
	var subProjectIndexes []pair[string, int]

	for i, h := range headerRow {
		if !strings.HasPrefix(h, projectName) {
			continue
		}

		parts := strings.Split(h, ":")
		switch len(parts) {
		case 1:
			mainProjectIndex = i
		case 2:
			subProjectIndexes = append(subProjectIndexes, pair[string, int]{parts[1], i})
		default:
			return nil
		}
	}

	if mainProjectIndex == -1 {
		return nil
	}

	return &columnIndexes{
		proj:     mainProjectIndex,
		subProjs: subProjectIndexes,
	}
}

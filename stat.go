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

package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
)

const (
	dateCol = "date"
)

func writeCSVUpdateIfNeeded(date string, properties []Properties) {

	fileName := writeFolderPath + "result.csv"
	tmpFileName := writeFolderPath + "result_tmp.csv"

	var newFile bool
	if _, err := os.Stat(fileName); err == nil {
		newFile = false
	} else if errors.Is(err, os.ErrNotExist) {
		newFile = true
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		// But ok here
	}

	var fileContent [][]string
	var file *os.File
	var originalFile *os.File

	if newFile {
		var err error
		file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		stopIfErrf("create csv file error: %w", err)
	} else {
		var err error
		originalFile, err = os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0644)
		stopIfErrf("open originalfile err: %w", err)
		file, err = os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		stopIfErrf("open csv file error: %w", err)
		reader := csv.NewReader(originalFile)
		fileContent, err = reader.ReadAll()
		stopIfErrf("read csv line: %w", err)

		originalFile.Close()
	}

	newContent := updateOriginal(date, hoursPerTomato, properties, fileContent)

	writer := csv.NewWriter(file)
	writer.WriteAll(newContent)
	writer.Flush()
	file.Close()

	if !newFile {
		stopIfErrf("remove originalFile: %w", os.Remove(fileName))
		stopIfErrf("rename tmp: %w", os.Rename(tmpFileName, fileName))
	}
}

func sortedKeys(m map[string]int) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}

	sort.Slice(ks, func(i, j int) bool {
		return ks[i] < ks[j]
	})
	return ks
}

// returns report sorted keys from properties
// returns map<key,values>
func reportKeys(props []Properties) ([]string, map[string]int) {
	var keys []string
	kToProps := map[string]int{}

	for _, p := range props {
		keys = append(keys, p.ReportKey)
		kToProps[p.ReportKey] = p.DoneTotal

		for _, k := range sortedKeys(p.Done) {
			composedKey := p.ReportKey + ":" + k
			keys = append(keys, composedKey)
			kToProps[composedKey] = p.Done[k]
		}
	}

	return keys, kToProps
}

// finds existing key indexes
func existingKeys(headerRow []string) map[string]int {
	keys := map[string]int{}
	for i, r := range headerRow {
		if i == 0 && r != dateCol {
			panic(fmt.Sprintf("not date in the first column: %v", headerRow))
		}

		keys[r] = i
	}
	return keys
}

func updateOriginal(date string, hoursPerTomato float64, properties []Properties, original [][]string) [][]string {
	if len(original) == 0 {
		original = [][]string{
			{dateCol},
		}
	}

	existingKeys := existingKeys(original[0])
	keysToReport, dataToReport := reportKeys(properties)

	totalColumns := len(existingKeys)
	for _, k := range keysToReport {
		if _, ok := existingKeys[k]; !ok {
			totalColumns++
		}
	}

	newReportRow := make([]string, totalColumns)
	for i := range newReportRow {
		newReportRow[i] = "0.0"
	}
	newReportRow[0] = date

	for _, k := range keysToReport {
		index, ok := existingKeys[k]
		if !ok {
			index = len(existingKeys)
			existingKeys[k] = index
			original[0] = append(original[0], k)
		}

		newReportRow[index] = strconv.FormatFloat(float64(dataToReport[k])*hoursPerTomato, 'f', 1, 64)
	}

	// looking for an already existing raw
	for i := len(original) - 1; i >= 0; i-- {
		if original[i][0] == date {
			original[i] = newReportRow

			return original
		}
	}

	return append(original, newReportRow)
}

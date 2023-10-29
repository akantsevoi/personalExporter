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
	var newFile bool
	if _, err := os.Stat("result.csv"); err == nil {
		newFile = false
	} else if errors.Is(err, os.ErrNotExist) {
		newFile = true
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		// But ok here
	}

	file, err := os.OpenFile("result.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	stopIfErrf("open csv file error: %w", err)
	defer file.Close()

	var fileContent [][]string
	if !newFile {
		reader := csv.NewReader(file)
		fileContent, err = reader.ReadAll()
		stopIfErrf("read csv line:%v %w", err, err.Error())
	}

	newContent := updateOriginal(date, properties, fileContent)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.WriteAll(newContent)
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

func updateOriginal(date string, properties []Properties, original [][]string) [][]string {
	if len(original) == 0 {
		csvData := [][]string{
			{dateCol},
			{date},
		}
		for _, r := range properties {
			csvData[0] = append(csvData[0], r.ReportKey)
			csvData[1] = append(csvData[1], strconv.Itoa(r.DoneTotal))

			for _, k := range sortedKeys(r.Done) {
				v := r.Done[k]
				csvData[0] = append(csvData[0], r.ReportKey+":"+k)
				csvData[1] = append(csvData[1], strconv.Itoa(v))
			}
		}
		return csvData
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
		newReportRow[i] = "0"
	}
	newReportRow[0] = date

	for _, k := range keysToReport {
		index, ok := existingKeys[k]
		if !ok {
			index = len(existingKeys)
			existingKeys[k] = index
			original[0] = append(original[0], k)
		}
		newReportRow[index] = strconv.Itoa(dataToReport[k])
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

// func createReportRow(date string, properties []Properties) []string {
// 	return nil
// }

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// detect that it's done but for subproject
	suffix = "_sub_done_today"
)

func mapToProperties(input map[string]string) *Properties {
	pr := Properties{
		Done: make(map[string]int),
	}

	for k, v := range input {
		name, hasSuffix := strings.CutSuffix(k, suffix)
		if hasSuffix {
			increment := toInt(v)
			pr.DoneTotal += increment
			pr.Done[name] = increment
		} else if k == "doneToday" {
			pr.DoneTotal = toInt(v)
		} else if k == "reportKey" {
			pr.ReportKey = input["reportKey"]
		} else if k == "priority" {
			pr.Priority = toInt(input["priority"])
		} else if k == "status" {
			pr.Status = input["status"]
		} else if k == "tomatosPerDay" {
			pr.TomatoesPerDay = toInt(input["tomatosPerDay"])
		} else {
			pr.OtherProps = append(pr.OtherProps, pair[string, string]{k, v})
		}
	}

	if pr.DoneTotal == 0 {
		return nil
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

func toFloat64(v string) float64 {
	i, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic(fmt.Errorf("parseFloat64: %v: %w", v, err))
	}
	return i
}

func shouldExist(path string, m map[string]string, key string) {
	if _, ok := m[key]; !ok {
		fmt.Printf("no %v in %v\n%v", key, path, m)
		os.Exit(1)
	}
}

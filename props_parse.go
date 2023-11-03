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
			if increment == 0 {
				continue
			}
			pr.DoneTotal += increment
			pr.Done[name] = increment

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

func shouldExist(path string, m map[string]string, key string) {
	if _, ok := m[key]; !ok {
		fmt.Printf("no %v in %v\n%v", key, path, m)
		os.Exit(1)
	}
}

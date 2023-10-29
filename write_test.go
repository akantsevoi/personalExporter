package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateOriginal(t *testing.T) {
	for _, c := range []struct {
		label      string
		date       string
		fromFile   [][]string
		properties []Properties
		exp        [][]string
	}{
		{
			label: "non_empty_file_already_reported",
			date:  "2043-12-04",
			properties: []Properties{
				{ReportKey: "key1", DoneTotal: 4},
			},
			fromFile: [][]string{
				{"date", "key1"},
				{"2043-12-04", "4"},
			},
			exp: [][]string{
				{"date", "key1"},
				{"2043-12-04", "4"},
			},
		},
		{
			label: "empty_file_new_data",
			date:  "2043-12-04",
			properties: []Properties{
				{ReportKey: "key1", DoneTotal: 4},
				{ReportKey: "key3", DoneTotal: 2, Done: map[string]int{
					"sub1": 1,
					"sub2": 0,
					"sub3": 1,
				}},
			},
			exp: [][]string{
				{"date", "key1", "key3", "key3:sub1", "key3:sub2", "key3:sub3"},
				{"2043-12-04", "4", "2", "1", "0", "1"},
			},
		},
	} {
		t.Run(c.label, func(t *testing.T) {
			res := updateOriginal(c.date, c.properties, c.fromFile)
			require.Equal(t, c.exp, res)
		})
	}
}

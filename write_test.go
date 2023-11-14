package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateOriginal(t *testing.T) {
	for _, c := range []struct {
		label          string
		date           string
		hoursPerTomato float64
		fromFile       [][]string
		properties     []Properties
		exp            [][]string
	}{
		{
			label:          "non_empty_file_new_row",
			date:           "2043-12-04",
			hoursPerTomato: 1.0,
			properties: []Properties{
				{ReportKey: "key3", DoneTotal: 4, Done: map[string]int{
					"sub1": 3,
					"sub2": 0,
					"sub3": 1,
				}},
				{ReportKey: "key2", DoneTotal: 2, Done: map[string]int{
					"sub1": 2,
					"sub2": 0,
				}},
			},
			fromFile: [][]string{
				{"date", "key1", "key1:sub1"},
				{"2043-12-03", "4", "4"},
			},
			exp: [][]string{
				{"date", "key1", "key1:sub1", "key3", "key3:sub1", "key3:sub3", "key2", "key2:sub1"},
				{"2043-12-03", "4", "4", "0.0", "0.0", "0.0", "0.0", "0.0"},
				{"2043-12-04", "0.0", "0.0", "4.0", "3.0", "1.0", "2.0", "2.0"},
			},
		},
		{
			label:          "non_empty_file_already_reported_differentKeys_with_subkeys",
			date:           "2043-12-04",
			hoursPerTomato: 1.0,
			properties: []Properties{
				{ReportKey: "key2", DoneTotal: 2, Done: map[string]int{
					"sub1": 1,
					"sub2": 1,
				}},
			},
			fromFile: [][]string{
				{"date", "key1"},
				{"2043-12-04", "4"},
			},
			exp: [][]string{
				{"date", "key1", "key2", "key2:sub1", "key2:sub2"},
				{"2043-12-04", "0.0", "2.0", "1.0", "1.0"},
			},
		},
		{
			label:          "non_empty_file_already_reported_differentKeys",
			date:           "2043-12-04",
			hoursPerTomato: 1.0,
			properties: []Properties{
				{ReportKey: "key2", DoneTotal: 2},
			},
			fromFile: [][]string{
				{"date", "key1"},
				{"2043-12-04", "4"},
			},
			exp: [][]string{
				{"date", "key1", "key2"},
				{"2043-12-04", "0.0", "2.0"},
			},
		},
		{
			label:          "non_empty_file_already_reported_with_updates",
			date:           "2043-12-04",
			hoursPerTomato: 1.0,
			properties: []Properties{
				{ReportKey: "key1", DoneTotal: 4},
				{ReportKey: "key2", DoneTotal: 2},
			},
			fromFile: [][]string{
				{"date", "key1"},
				{"2043-12-04", "4"},
			},
			exp: [][]string{
				{"date", "key1", "key2"},
				{"2043-12-04", "4.0", "2.0"},
			},
		},
		{
			label:          "non_empty_file_already_reported",
			date:           "2043-12-04",
			hoursPerTomato: 1.0,
			properties: []Properties{
				{ReportKey: "key1", DoneTotal: 4},
			},
			fromFile: [][]string{
				{"date", "key1"},
				{"2043-12-04", "4.0"},
			},
			exp: [][]string{
				{"date", "key1"},
				{"2043-12-04", "4.0"},
			},
		},
		{
			label:          "empty_file_new_data_not_round_hours",
			date:           "2043-12-04",
			hoursPerTomato: 0.3,
			properties: []Properties{
				{ReportKey: "key1", DoneTotal: 4},
				{ReportKey: "key3", DoneTotal: 2, Done: map[string]int{
					"sub1": 1,
					"sub2": 0,
					"sub3": 1,
				}},
			},
			exp: [][]string{
				{"date", "key1", "key3", "key3:sub1", "key3:sub3"},
				{"2043-12-04", "1.2", "0.6", "0.3", "0.3"},
			},
		},
		{
			label:          "empty_file_new_data",
			date:           "2043-12-04",
			hoursPerTomato: 1,
			properties: []Properties{
				{ReportKey: "key1", DoneTotal: 4},
				{ReportKey: "key3", DoneTotal: 2, Done: map[string]int{
					"sub1": 1,
					"sub2": 0,
					"sub3": 1,
				}},
			},
			exp: [][]string{
				{"date", "key1", "key3", "key3:sub1", "key3:sub3"},
				{"2043-12-04", "4.0", "2.0", "1.0", "1.0"},
			},
		},
	} {
		t.Run(c.label, func(t *testing.T) {
			res := updateOriginal(c.date, c.hoursPerTomato, c.properties, c.fromFile)
			require.Equal(t, c.exp, res)
		})
	}
}

func TestWriteCSVUpdate(t *testing.T) {
	testFolderPath := "tmp_test_path/"

	hoursPerTomato = 0.5

	err := os.RemoveAll(testFolderPath)
	require.NoError(t, err)

	err = os.Mkdir(testFolderPath, os.FileMode(0755))
	require.NoError(t, err)

	t.Run("create_new_result", func(t *testing.T) {
		// basic_export.csv.golden
		writeCSVUpdateIfNeeded("2023-04-22", testFolderPath, []Properties{
			{ReportKey: "repKey1", DoneTotal: 10},
			{ReportKey: "repKey2", DoneTotal: 1},
		})

		result := helperFileContent(t, testFolderPath+"result.csv")
		want := helperFileContent(t, "test_files/basic_export.csv.golden")
		require.Equal(t, want, result)
	})

	t.Run("add_to_existed", func(t *testing.T) {
		copyFile(t, "test_files/existedResult.csv", testFolderPath+"result.csv")
		writeCSVUpdateIfNeeded("2023-06-29", testFolderPath, []Properties{
			{ReportKey: "proj1", DoneTotal: 3, Done: map[string]int{
				"subproj":  2,
				"subproj2": 1,
			}},
		})

		result := helperFileContent(t, testFolderPath+"result.csv")
		want := helperFileContent(t, "test_files/after_export_to_existed.csv.golden")
		require.Equal(t, want, result)
	})

	t.Run("clean_up", func(t *testing.T) {
		err := os.RemoveAll(testFolderPath)
		require.NoError(t, err)
	})
}

func copyFile(t *testing.T, from, to string) {
	srcFile, err := os.Open(from)
	require.NoError(t, err)
	defer srcFile.Close()

	destFile, err := os.Create(to)
	require.NoError(t, err)
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	require.NoError(t, err)

	err = destFile.Sync()
	require.NoError(t, err)
}

func helperFileContent(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Error loading golden file: %s", err)
	}
	return string(content)
}

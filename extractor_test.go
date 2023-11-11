package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractPropertiesFromFile(t *testing.T) {
	for _, c := range []struct {
		label    string
		filePath string
		expProps *Properties
		expErr   *string
	}{
		{
			label:    "project_with_subprojects",
			filePath: "test_files/projectWithSubprojects.md",
			expProps: &Properties{
				DoneTotal:      2,
				ReportKey:      "repKeySub",
				Status:         "hold",
				Priority:       2,
				TomatoesPerDay: 8,
				Done: map[string]int{
					"myProj1": 1,
					"myProj2": 1,
					"other":   0,
				},
			},
		},
		{
			label:    "project_file_props",
			filePath: "test_files/correctProjectFile.md",
			expProps: &Properties{
				DoneTotal:      2,
				ReportKey:      "repKey1",
				Status:         "progress",
				Priority:       3,
				TomatoesPerDay: 2,
				Done:           make(map[string]int),
			},
		},
		{
			label:    "not_project_file_no_props",
			filePath: "test_files/notProjectFile.md",
		},
		{
			label:    "not_exists_file",
			filePath: "test_files/notExists.md",
			expErr:   ptr("no such file or directory"),
		},
	} {
		t.Run(c.label, func(t *testing.T) {
			props, err := extractFrontmatterProperties(c.filePath)
			if c.expErr == nil {
				require.NoError(t, err)
				require.Equal(t, c.expProps, props)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), *c.expErr)
			}
		})
	}
}

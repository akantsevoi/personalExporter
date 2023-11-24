package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResetProgress(t *testing.T) {
	testFolderPath := "tmp_test_path/"

	err := os.RemoveAll(testFolderPath)
	require.NoError(t, err)

	err = os.Mkdir(testFolderPath, os.FileMode(0755))
	require.NoError(t, err)

	// resetProgress()

	t.Run("reset_single_project_file", func(t *testing.T) {
		copyFile(t, "test_files/correctProjectFile.md", testFolderPath+"correctProjectFile.md")
		err := resetProgress([]string{
			testFolderPath + "correctProjectFile.md",
		})
		require.NoError(t, err)

		result := helperFileContent(t, testFolderPath+"correctProjectFile.md")
		expected := helperFileContent(t, "test_files/correctProjectFile_reseted.md.golden")
		require.Equal(t, expected, result)
	})

	t.Run("reset_subproject_file", func(t *testing.T) {
		copyFile(t, "test_files/projectWithSubprojects.md", testFolderPath+"projectWithSubprojects.md")
		err := resetProgress([]string{
			testFolderPath + "projectWithSubprojects.md",
		})
		require.NoError(t, err)

		result := helperFileContent(t, testFolderPath+"projectWithSubprojects.md")
		expected := helperFileContent(t, "test_files/projectWithSubprojects_reseted.md.golden")
		require.Equal(t, expected, result)
	})

	t.Run("reset_subproject_file_with_some_extra_props", func(t *testing.T) {
		copyFile(t, "test_files/projectWithSubprojects_extra_props.md", testFolderPath+"projectWithSubprojects_extra_props.md")
		err := resetProgress([]string{
			testFolderPath + "projectWithSubprojects_extra_props.md",
		})
		require.NoError(t, err)

		result := helperFileContent(t, testFolderPath+"projectWithSubprojects_extra_props.md")
		expected := helperFileContent(t, "test_files/projectWithSubprojects_reseted_extra_props.md.golden")
		require.Equal(t, expected, result)
	})

	t.Run("clean_up", func(t *testing.T) {
		err := os.RemoveAll(testFolderPath)
		require.NoError(t, err)
	})
}

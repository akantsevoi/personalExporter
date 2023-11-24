package main

import (
	"fmt"
	"os"
	"strings"
)

func resetProgress(filePaths []string) error {
	for _, p := range filePaths {
		props, err := extractFrontmatterProperties(p)
		if err != nil {
			return fmt.Errorf("resetProgress: error extracting %s: %w", p, err)
		}
		if props == nil {
			continue
		}

		content, err := fileContent(p)
		if err != nil {
			return fmt.Errorf("resetProgress: getContent %s: %w", p, err)
		}
		lines := strings.SplitN(content, "---", 3)

		lines[1] = resetedStringProps(props)

		newContent := strings.Join(lines, "---")

		if err := os.WriteFile(p, []byte(newContent), 0666); err != nil {
			return fmt.Errorf("resetProgress: writeNewContent %s: %w", p, err)
		}
	}
	return nil
}

func fileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("fileContent: %w", err)
	}
	return string(content), nil
}

func resetedStringProps(props *Properties) string {
	res := fmt.Sprintf(`
status: %v
tomatosPerDay: %v
reportKey: %v
priority: %v
`,
		props.Status, props.TomatoesPerDay, props.ReportKey, props.Priority)

	if len(props.Done) == 0 {
		res += "doneToday: 0\n"
	} else {
		for _, k := range sortedKeys(props.Done) {
			res += fmt.Sprintf("%s_sub_done_today: 0\n", k)
		}
	}

	for _, p := range props.OtherProps {
		res += fmt.Sprintf("%s: %s\n", p.First, p.Second)
	}

	return res
}

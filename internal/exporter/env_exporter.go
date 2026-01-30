package exporter

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

type EnvExporter struct{}

func (e *EnvExporter) Format() string {
	return "env"
}

func (e *EnvExporter) Export(secrets map[string]string, opts ExportOptions) error {
	if !opts.Append && !opts.Force && fileExists(opts.TargetPath) {
		return fmt.Errorf("file %s already exists (use --force to overwrite or --append to add to it)", opts.TargetPath)
	}

	preview, err := e.Preview(secrets, opts)
	if err != nil {
		return err
	}

	if opts.DryRun {
		return nil
	}

	if opts.Append && fileExists(opts.TargetPath) {
		return e.appendToFile(secrets, preview, opts)
	}

	return e.writeNewFile(preview, opts)
}

func (e *EnvExporter) Preview(secrets map[string]string, opts ExportOptions) (*Preview, error) {
	preview := &Preview{
		NewKeys:     make([]string, 0),
		UpdatedKeys: make([]string, 0),
		SkippedKeys: make([]string, 0),
	}

	existingKeys := make(map[string]string)
	if opts.Append && fileExists(opts.TargetPath) {
		existingKeys = e.parseExistingFile(opts.TargetPath)
	}

	var content strings.Builder
	sortedKeys := sortKeys(secrets)

	if opts.Append && fileExists(opts.TargetPath) {
		data, err := os.ReadFile(opts.TargetPath)
		if err != nil {
			return nil, err
		}
		content.Write(data)
		if !strings.HasSuffix(content.String(), "\n") {
			content.WriteString("\n")
		}
		content.WriteString(fmt.Sprintf("\n# Added by veil on %s\n", time.Now().Format("2006-01-02T15:04:05Z")))
	}

	for _, key := range sortedKeys {
		value := secrets[key]

		if existingValue, exists := existingKeys[key]; exists {
			if existingValue == value && !opts.Force {
				preview.SkippedKeys = append(preview.SkippedKeys, key)
				continue
			}
			if opts.Force {
				preview.UpdatedKeys = append(preview.UpdatedKeys, key)
			} else {
				preview.SkippedKeys = append(preview.SkippedKeys, key)
				continue
			}
		} else {
			preview.NewKeys = append(preview.NewKeys, key)
		}

		content.WriteString(fmt.Sprintf("%s=%s\n", key, e.escapeValue(value)))
	}

	preview.Content = content.String()
	return preview, nil
}

func (e *EnvExporter) writeNewFile(preview *Preview, opts ExportOptions) error {
	return safeWriteFile(opts.TargetPath, []byte(preview.Content), 0600, opts.Backup, opts.BackupDir)
}

func (e *EnvExporter) appendToFile(secrets map[string]string, preview *Preview, opts ExportOptions) error {
	if len(preview.NewKeys) == 0 && len(preview.UpdatedKeys) == 0 {
		return nil
	}

	existingContent, err := os.ReadFile(opts.TargetPath)
	if err != nil {
		return err
	}

	var newContent strings.Builder
	newContent.Write(existingContent)

	if !strings.HasSuffix(newContent.String(), "\n") {
		newContent.WriteString("\n")
	}

	newContent.WriteString(fmt.Sprintf("\n# Added by veil on %s\n", time.Now().Format("2006-01-02T15:04:05Z")))

	for _, key := range preview.NewKeys {
		newContent.WriteString(fmt.Sprintf("%s=%s\n", key, e.escapeValue(secrets[key])))
	}

	if opts.Force {
		existingKeys := e.parseExistingFile(opts.TargetPath)
		for _, key := range preview.UpdatedKeys {
			existingKeys[key] = secrets[key]
		}

		var updatedContent strings.Builder
		scanner := bufio.NewScanner(bytes.NewReader(existingContent))
		for scanner.Scan() {
			line := scanner.Text()
			if key, _, found := strings.Cut(line, "="); found {
				if newVal, exists := existingKeys[key]; exists {
					updatedContent.WriteString(fmt.Sprintf("%s=%s\n", key, e.escapeValue(newVal)))
					delete(existingKeys, key)
					continue
				}
			}
			updatedContent.WriteString(line)
			updatedContent.WriteString("\n")
		}
		newContent.Reset()
		newContent.WriteString(updatedContent.String())
	}

	return safeWriteFile(opts.TargetPath, []byte(newContent.String()), 0600, opts.Backup, opts.BackupDir)
}

func (e *EnvExporter) parseExistingFile(path string) map[string]string {
	result := make(map[string]string)
	data, err := os.ReadFile(path)
	if err != nil {
		return result
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if key, value, found := strings.Cut(line, "="); found {
			result[key] = e.unescapeValue(value)
		}
	}

	return result
}

func (e *EnvExporter) escapeValue(value string) string {
	if strings.Contains(value, " ") || strings.Contains(value, "\t") {
		return fmt.Sprintf("\"%s\"", strings.ReplaceAll(value, "\"", "\\\""))
	}
	return value
}

func (e *EnvExporter) unescapeValue(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = value[1 : len(value)-1]
		value = strings.ReplaceAll(value, "\\\"", "\"")
	}
	return value
}

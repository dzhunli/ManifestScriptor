package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

func validateYAML(yamlContent string) error {
	var data interface{}
	err := yaml.Unmarshal([]byte(yamlContent), &data)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	templatePath := flag.String("t", "", "Path to the YAML template file")
	scriptPath := flag.String("s", "", "Path to the Bash script file")
	outputPath := flag.String("o", "output.yaml", "Path to the output YAML file")
	flag.Parse()

	if *templatePath == "" || *scriptPath == "" {
		color.Red("Error: Both -t and -s arguments are required.")
		flag.Usage()
		os.Exit(1)
	}

	templateContent, err := os.ReadFile(*templatePath)
	if err != nil {
		color.Red("✖ ERROR: Failed to read template file: %v", err)
		os.Exit(1)
	}

	scriptContent, err := os.ReadFile(*scriptPath)
	if err != nil {
		color.Red("✖ ERROR: Failed to read script file: %v", err)
		os.Exit(1)
	}

	finalOutput, err := replaceScriptWithIndentation(string(templateContent), string(scriptContent))
	if err != nil {
		color.Red("✖ ERROR: %v", err)
		os.Exit(1)
	}

	err = os.WriteFile(*outputPath, []byte(finalOutput), 0644)
	if err != nil {
		color.Red("✖ ERROR: Failed to write output file: %v", err)
		os.Exit(1)
	}

	color.Cyan("Validating YAML...")
	err = validateYAML(finalOutput)
	if err != nil {
		color.Red("✖ ERROR: YAML validation failed: %v", err)
		color.Yellow("------ Generated YAML ------")
		color.Yellow(finalOutput)
		os.Exit(1)
	}

	color.Green("✔ SUCCESS: YAML is valid and written to '%s'", *outputPath)
}

func replaceScriptWithIndentation(templateContent, scriptContent string) (string, error) {
	templateLines := strings.Split(templateContent, "\n")
	var lineIndex int
	var indentLevel int

	for i, line := range templateLines {
		if strings.Contains(line, "{|script|}") {
			lineIndex = i
			indentLevel = countLeadingSpaces(line)
			break
		}
	}

	if lineIndex == 0 {
		return "", errors.New("{|script|} placeholder not found in the template")
	}

	formattedScript := formatScript(scriptContent, indentLevel) // Отступ на том же уровне, что и {script}

	// Заменяем строку с {script} на весь скрипт с правильным отступом
	templateLines[lineIndex] = formattedScript

	return strings.Join(templateLines, "\n"), nil
}

func formatScript(script string, baseIndent int) string {
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(script))
	indent := strings.Repeat(" ", baseIndent)

	for scanner.Scan() {
		line := scanner.Text()
		buffer.WriteString(indent + line + "\n")
	}

	if err := scanner.Err(); err != nil {
		color.Red("✖ ERROR: Failed to read script: %v", err)
		os.Exit(1)
	}

	return buffer.String()
}

func countLeadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

package main

import (
	"bufio"
	"bytes"
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

	formattedScript := formatScript(string(scriptContent))
	finalOutput := strings.Replace(string(templateContent), "{|script|}", formattedScript, 1)

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

func formatScript(script string) string {
	var buffer bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(script))
	for scanner.Scan() {
		line := scanner.Text()
		buffer.WriteString("            " + strings.ReplaceAll(line, "\t", "    ") + "\n")
	}
	if err := scanner.Err(); err != nil {
		color.Red("✖ ERROR: Failed to read script: %v", err)
		os.Exit(1)
	}
	return buffer.String()
}

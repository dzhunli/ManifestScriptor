# From Template and Script ==> YAML

![Build and Release](https://github.com/dzhunli/yaml-mscr/actions/workflows/realese.yml/badge.svg)
![Contributions](https://img.shields.io/badge/contributions-welcome-brightgreen)
![Go Version](https://img.shields.io/github/go-mod/go-version/dzhunli/yaml-mscr?filename=app%2Fgo.mod)

## Overview

This project is a Go-based tool for generating YAML files from a predefined template and a Bash script. 

The main goal is to replace the `{|script|}` marker in a YAML template with the formatted content of a Bash script. The resulting YAML file is then validated for correctness. 

### Features
- Generates YAML files by merging templates with Bash script content.
- Automatically validates the generated YAML.
- CLI with customizable options for file paths and output.
- User-friendly error messages for easier debugging.
- Simplifies the integration of Bash scripts into YAML templates.

## Usage

### How to Use

Download the pre-built binary for your platform from the [Releases page](https://github.com/dzhunli/yaml-mscr/releases).

Run the tool using the following command:

```bash
./manscr -t path/to/template.yaml -s path/to/script.sh -o path/to/output.yaml


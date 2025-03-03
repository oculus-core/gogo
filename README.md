# Gogo CLI

[![Release](https://github.com/oculus-core/gogo/actions/workflows/release.yml/badge.svg)](https://github.com/oculus-core/gogo/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)](https://golang.org/doc/go1.24)

CLI tool for generating Go projects with best practices.

## Overview

Gogo scaffolds Go projects based on a reference workspace. It includes:

- Interactive project setup wizard
- Project generation with customizable templates
- Homebrew formula support
- Commit message validation using Conventional Commits format
- Automated development tooling setup (assumes VSCode)
- Support for different project types (CLI, API, Library)
- Configuration file support

## TODO

The following features are planned for future releases:

- Integration with linters and pre-commit hooks
- GitHub Actions workflows for CI/CD
- Observability tooling with OpenTelemetry
- Quality Go frameworks integration (Cobra, Viper, Gin, etc.)

---

## Installation

### Prerequisites

- Go 1.24 or later
- Git

### Installation Options

#### Using Go Install

```bash
go install github.com/oculus-core/gogo@latest
```

#### Download Binary

Download from the [releases page](https://github.com/oculus-core/gogo/releases/latest).

#### Package Managers

##### Homebrew (macOS/Linux)

```bash
# Latest version
brew install oculus-core/homebrew-gogo/gogo

# Specific version series
brew install oculus-core/homebrew-gogo/gogo@0.1
```

Available at [oculus-core/homebrew-gogo](https://github.com/oculus-core/homebrew-gogo).

### Build from Source

```bash
git clone https://github.com/oculus-core/gogo.git
cd gogo
go build -o bin/gogo
```

---

## Usage

```bash
# Create a new project with wizard
gogo new my-project

# Create project in specific directory
gogo new my-project --output /path/to/output

# Create project with default settings
gogo new my-project --skip-wizard

# Create project with specific project type
gogo new my-project --type cli
gogo new my-project --type api
gogo new my-project --type library

# Create project from configuration file
gogo new my-project --config path/to/config.yaml

# Show version
gogo version

# Show help
gogo help
```

## Project Types

Gogo supports different project types, each with its own structure and dependencies:

### CLI Applications

```bash
gogo new my-cli-app --type cli
```

- Command-line interface structure
- Includes Cobra for command handling
- Includes Viper for configuration management
- Includes version command

### API Applications

```bash
gogo new my-api --type api
```

- Web service structure
- Includes Gin web framework
- Includes configuration management
- Basic API endpoints (health check, hello world)

### Library/Package

```bash
gogo new my-lib --type library
```

- Package-oriented structure
- No cmd directory
- Includes test files
- Ready for distribution as a Go module

## Configuration File

You can use a YAML configuration file to define your project settings:

```yaml
# Gogo Project Configuration
name: my-awesome-project
module: github.com/username/my-awesome-project
description: A sample Go project created with Gogo
license: MIT
author: Your Name
type: cli  # Options: default, cli, api, library

# Project structure options
use_cmd: true
use_internal: true
use_pkg: true
use_test: true
use_docs: true
create_readme: true
create_license: true
create_makefile: true

# Code quality tools
use_linters: true
use_pre_commit_hooks: true
use_git_hooks: true

# Dependencies
use_cobra: true
use_viper: true
use_gin: false

# CI/CD
use_github_actions: true
```

Use the configuration file with:

```bash
gogo new my-project --config path/to/config.yaml
```

## Wizard Process

When running `gogo new my-project`, you'll go through:

1. **Project Information**
   - Project name
   - Module path
   - Description
   - Author
   - License
   - Project type

2. **Project Structure**
   - Select directories (cmd, internal, pkg, etc.)
   - Choose files to generate
   - Set up dependencies

3. **Project Generation**
   - Review selections
   - Create project

## Development

```bash
git clone https://github.com/oculus-core/gogo.git
cd gogo
go mod tidy
make build
make test
```

### Project Structure

- `cmd/`: Command-line interface
- `internal/`: Private application code
- `pkg/`: Public API packages
- `examples/`: Example configuration files

## Contributing

Contributions are welcome! Here's how you can contribute:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-new-feature`
3. Make your changes
4. Run tests: `make test`
5. Commit your changes: `git commit -am 'Add some feature'`
6. Push to the branch: `git push origin feature/my-new-feature`
7. Submit a pull request

Make sure your code follows the project's coding standards and includes appropriate tests.

## License

MIT

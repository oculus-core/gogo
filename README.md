# Gogo CLI 🚀

[![Release](https://github.com/oculus-core/gogo/actions/workflows/release.yml/badge.svg)](https://github.com/oculus-core/gogo/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8.svg)](https://golang.org/doc/go1.16)

A powerful and intuitive CLI tool for generating professional Go projects with best practices.
Accelerate your Go development with structured templates, tooling, and automation.

## Overview

Gogo helps scaffold Go projects based on a reference workspace. It includes:

- 🧙 Interactive project setup wizard
- 🏗️ Support for various project structures, based on project type (API, CLI, etc.)
- 🔍 Integration with linters and pre-commit hooks to ensure code quality
- 📝 Commit message validation using Conventional Commits format
- 🔄 GitHub Actions workflows for CI/CD
- 📋 Customizable templates
- 🛠️ Automated development tooling setup
- 📊 Observability tooling (metrics, tracing, logging) using OpenTelemetry out the box
- 📚 Use of quality Go frameworks and libraries (eg. Cobra, Viper, Gin, etc.)

---

## Installation

### Prerequisites

- Go 1.16 or later
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

# Show version
gogo version

# Show help
gogo help
```

## Wizard Process

When running `gogo new my-project`, you'll go through:

1. **Project Information**
   - Project name
   - Module path
   - Description
   - Author
   - License

2. **Project Structure**
   - Select directories (cmd, internal, pkg, etc.)
   - Choose files to generate
   - Configure code quality tools
   - Set up dependencies and CI/CD

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

## Contributing

Contributions are welcome! Here's how you can contribute:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-new-feature`
3. Make your changes
4. Run tests: `make test`
5. Commit your changes: `git commit -am 'Add some feature'`
6. Push to the branch: `git push origin feature/my-new-feature`
7. Submit a pull request

Please make sure your code follows the project's coding standards and includes appropriate tests.

## License

MIT

[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/orzkratos/supervisorkratos/release.yml?branch=main&label=BUILD)](https://github.com/orzkratos/supervisorkratos/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/orzkratos/supervisorkratos)](https://pkg.go.dev/github.com/orzkratos/supervisorkratos)
[![Coverage Status](https://img.shields.io/coveralls/github/orzkratos/supervisorkratos/main.svg)](https://coveralls.io/github/orzkratos/supervisorkratos?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/orzkratos/supervisorkratos.svg)](https://github.com/orzkratos/supervisorkratos/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/orzkratos/supervisorkratos)](https://goreportcard.com/report/github.com/orzkratos/supervisorkratos)

# supervisorkratos

Go package to generate supervisord configuration with Kratos microservices integration.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->
## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Main Features

🎯 **Fluent Configuration API**: Chain methods to build supervisord config
⚡ **Kratos Integration**: Optimized configuration patterns to run Kratos microservices
🔄 **Group Management**: Multi-program groups with centralized configuration
🌍 **Tested Configuration**: Battle-tested templates that run high-performance services
📋 **Strong Typing**: Typed configuration with sensible defaults

## Installation

```bash
go get github.com/orzkratos/supervisorkratos
```

## Usage

### Single Program Configuration

```go
package main

import (
    "fmt"
    "github.com/orzkratos/supervisorkratos"
)

func main() {
    // Create program config with required parameters
    program := supervisorkratos.NewProgramConfig(
        "myapp",           // Program name
        "/opt/myapp",      // Program root DIR
        "deploy",          // User name
        "/var/log/myapp",  // Log DIR
    ).WithStartRetries(10).
      WithEnvironment(map[string]string{
          "APP_ENV": "production",
      })

    // Generate configuration
    config := supervisorkratos.GenerateProgramConfig(program)
    fmt.Println(config)
}
```

### Group Configuration

```go
// Create multiple programs
apiServer := supervisorkratos.NewProgramConfig(
    "api-server", "/opt/api-server", "deploy", "/var/log/services",
).WithStartRetries(3)

worker := supervisorkratos.NewProgramConfig(
    "worker", "/opt/worker", "deploy", "/var/log/services",
).WithAutoStart(false)

// Create group
group := supervisorkratos.NewGroupConfig("microservices").
    AddProgram(apiServer).
    AddProgram(worker)

config := supervisorkratos.GenerateGroupConfig(group)
```

### Advanced Configuration

```go
// High-performance service configuration
program := supervisorkratos.NewProgramConfig(
    "high-perf", "/opt/high-perf", "performance", "/var/log/perf",
).WithStartRetries(100).
  WithStopWaitSecs(60).
  WithLogMaxBytes("500MB").
  WithLogBackups(50).
  WithPriority(1)
```

### Multi-Instance Deployment

```go
// Multi-instance web server
program := supervisorkratos.NewProgramConfig(
    "web-server", "/opt/web-server", "deploy", "/var/log/cluster",
).WithNumProcs(3).
  WithProcessName("%(program_name)s_%(process_num)02d").
  WithEnvironment(map[string]string{
      "PORT_BASE": "8080",
  })
```

## Configuration Options

### Process Settings
- `WithAutoStart(bool)` - Auto start on supervisord startup
- `WithAutoRestart(bool)` - Auto restart on crash
- `WithAutoRestartMode(string)` - Auto restart mode ("false"/"true"/"unexpected")
- `WithStartRetries(int)` - Max start attempts count
- `WithStartSecs(int)` - Wait time in seconds before start succeeds

### Logging
- `WithLogMaxBytes(string)` - Max log file size (e.g., "50MB", "1GB")
- `WithLogBackups(int)` - Log backup files count
- `WithRedirectStderr(bool)` - Redirect stderr to stdout

### Process Execution
- `WithStopWaitSecs(int)` - Clean stop timeout seconds
- `WithStopSignal(string)` - Stop command name (TERM, INT, QUIT)
- `WithKillAsGroup(bool)` - Terminate child processes as group
- `WithPriority(int)` - Start rank (small ranks start first)

### Multi-Instance
- `WithNumProcs(int)` - Count of process instances
- `WithProcessName(string)` - Process name template

### Environment
- `WithEnvironment(map[string]string)` - Environment setting
- `WithExitCodes([]int)` - Expected exit codes

## Recommended Workflow

```bash
# 1. Generate config file
go run main.go > /etc/supervisord/conf.d/myapp.conf

# 2. Reload supervisord
sudo supervisorctl reread
sudo supervisorctl update

# 3. Manage services
sudo supervisorctl start myapp
sudo supervisorctl status
```

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a mistake?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes and use significant commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/orzkratos/supervisorkratos.svg?variant=adaptive)](https://starchart.cc/orzkratos/supervisorkratos)
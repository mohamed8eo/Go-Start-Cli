# GoStart

A CLI tool to bootstrap Go projects with production-ready structure and search/install packages without leaving your terminal.

---

## Description

GoStart is a developer-friendly CLI that eliminates the friction of starting new Go projects. Instead of manually creating directories, remembering package URLs, and copy-pasting import paths, GoStart gives you:

- **Interactive project scaffolding** with a beautiful TUI
- **Pre-configured project structure** following Go best practices
- **Built-in package search** to find and install dependencies without opening a browser

Whether you're starting a new API service, a worker, or any Go project, GoStart sets up everything in seconds.

---

## Why I Built This

As a Go developer, I kept running into the same friction points:

1. **The "Package URL Dance"** — I'd be in the flow, writing code, and need a package like `argon2` or `jwt`. I'd have to stop, open a browser, search Google for the GitHub URL, copy it, and then `go get` it. This constant context-switching broke my concentration.

2. **Boilerplate Fatigue** — Every new project meant creating the same directory structure: `cmd/api/`, `internal/`, `pkg/`, `config/`, etc. It was tedious and inconsistent.

GoStart solves both problems. Now I can:
- Search packages directly from the terminal with an interactive interface
- Scaffold a complete project structure in under 30 seconds
- Stay in my editor and maintain focus

---

## Features

- **Interactive TUI** powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Framework Support**: Gin, Echo, Fiber, Chi, or standard library
- **Database Integration**: PostgreSQL, MySQL, SQLite, MongoDB
- **SQL Drivers**: GORM, sqlx, sqlc, pgx
- **Package Search**: Search GitHub for Go packages and install them interactively
- **Production-Ready Structure**: Separates `cmd/api`, `cmd/worker`, `internal/`, `pkg/`, `config/`, and more

---

## Installation

### From Source

```bash
go install github.com/mohamed8eo/gostart@latest
```

### Clone and Build

```bash
git clone https://github.com/mohamed8eo/gostart.git
cd gostart
go build -o gostart
mv gostart $GOPATH/bin/
```

### Requirements

- Go 1.21 or later
- Git (for fetching packages)

---

## Usage

### Initialize a New Project

Launch the interactive TUI:

```bash
gostart init
```

Or use CLI flags:

```bash
gostart init my-project --framework Gin --database PostgreSQL --sql GORM
```

**Available Flags:**
| Flag | Options |
|------|---------|
| `--framework` / `-f` | `Gin`, `Echo`, `Fiber`, `Chi`, `None` |
| `--database` / `-d` | `PostgreSQL`, `MySQL`, `SQLite`, `MongoDB`, `None` |
| `--sql` / `-s` | `GORM`, `sqlx`, `sqlc`, `pgx`, `None` |

**Created Structure:**
```
my-project/
├── cmd/
│   ├── api/
│   │   └── main.go          # API server entry point
│   └── worker/
│       └── main.go          # Background worker entry point
├── internal/
│   ├── auth/                # Authentication logic
│   ├── storage/             # Database/storage layer
│   └── transport/           # HTTP handlers/transports
├── pkg/
│   ├── logger/              # Shared logging utilities
│   └── crypto/              # Cryptographic helpers
├── api/
│   └── openapi.yaml         # API specification
├── config/
│   └── config.yaml          # Configuration files
├── scripts/
│   └── deploy.sh            # Deployment scripts
├── go.mod
└── README.md
```

### Search and Add Packages

Search for Go packages on GitHub and install them without leaving your terminal:

```bash
gostart add
```

This launches an interactive search where you can:
- Type to search packages (e.g., `jwt`, `argon2`, `redis`)
- Browse results with star counts
- Select packages to install with a single keystroke
- Add multiple packages to a queue and install them all at once

**Pro Tip:** Set a `GITHUB_TOKEN` environment variable for higher API rate limits:

```bash
export GITHUB_TOKEN=your_token_here
gostart add
```

---

## Examples

### Create a Gin API with PostgreSQL

```bash
gostart init my-api --framework Gin --database PostgreSQL --sql GORM
```

### Create a minimal worker project

```bash
gostart init my-worker --framework None --database None --sql None
```

### Search and install packages
![[Pasted image 20260405164132.png]]
---

## Contributing

Contributions are welcome! Here's how to get started:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/my-feature`
3. **Make your changes** and add tests if applicable
4. **Run the tests**: `go test ./...`
5. **Commit your changes**: `git commit -am 'Add my feature'`
6. **Push to your branch**: `git push origin feature/my-feature`
7. **Open a Pull Request**

### Development Setup

```bash
git clone https://github.com/mohamed8eo/gostart.git
cd gostart
go mod download
go build .
```

### Areas for Contribution

- Add support for more frameworks (e.g., Beego, Revel)
- Add Docker initialization option
- Add CI/CD template generation
- Improve search with pkg.go.dev integration
- Add tests for TUI components

### Code Style

- Follow standard Go conventions (`gofmt`, `golint`)
- Keep functions focused and well-documented
- Ensure the TUI remains responsive and keyboard-friendly

---

## License

[MIT](LICENSE)

---

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — Styling

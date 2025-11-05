## dev-env

A small Go CLI for managing local development projects (start/stop/list) using a simple JSON config in your home directory.

### Features

- Initialize a workspace at $HOME/dev-env
- Track projects in projects.json
- Start commands as detached processes and store their PIDs in stats.json
- Stop running projects by name
- List projects and show RUNNING/STOPPED state

### Requirements

- Go 1.20+
- Windows is the primary target (detached process flags in `start` use Windows APIs). Other platforms may require adjustments.

### Install

```bash
go build -o dev-env
# Optionally install to GOPATH/bin
go install ./...
```

Add the binary to your PATH if needed.

### Quick Start

```bash
# 1) Initialize workspace under your home directory
dev-env init

# 2) Add projects to $HOME/dev-env/projects.json (see example below; path is optional and the current project directory will be used if not provided)
dev-env add --name --cmd --path

# 3) List projects (the name is optional)
dev-env list --name

# 4) Start a project
dev-env start <project-name>

# 5) Stop a project
dev-env stop <project-name>
```

### Data files

- `$HOME/dev-env/projects.json`: Array of projects. Expected fields:

  - `Name`: Unique project name
  - `Path`: Working directory for the command
  - `Cmd`: Shell command to start the project (string, split by spaces)

- `$HOME/dev-env/stats.json`: Map of `projectName -> PID` for running projects.

Example `projects.json`:

```json
[
  {
    "Name": "web",
    "Path": "D:/Work/web-app",
    "Cmd": "npm run dev"
  },
  {
    "Name": "api",
    "Path": "D:/Work/api",
    "Cmd": "go run main.go"
  }
]
```

---

### Commands

#### Root command

- File: `root.go`
- Type: `*cobra.Command` (root)
- Description: Sets up the CLI root (`dev-env`) and wires subcommands.

```bash
dev-env --help
```

#### init

- File: `init.go`
- Type: `*cobra.Command` (subcommand)
- Use: `dev-env init`
- Short: Initialize the tool for development
- What it does:
  - Creates `$HOME/dev-env/`
  - Writes empty `projects.json` ([]) and `stats.json` ({}) files
- Example:

```bash
dev-env init
Initialized dev-env in C:/Users/<you>/dev-env
```

#### list

- File: `list.go`
- Type: `*cobra.Command` (subcommand)
- Use: `dev-env list`
- Flags:
  - `--name <project-name>`: Show details for a single project (with path and state)
- Behavior:
  - Reads `projects.json` and `stats.json`
  - Prints `[RUNNING]` or `[STOPPED]` per project
- Examples:

```bash
# List all
dev-env list
[RUNNING] - web
[STOPPED] - api

# Show one by name
dev-env list --name web
[RUNNING] - web - D:/Work/web-app
```

#### start

- File: `start.go`
- Type: `*cobra.Command` (subcommand)
- Use: `dev-env start <project-name>`
- Behavior:
  - Looks up the project by name in `projects.json`
  - Splits `Cmd` by spaces and starts it in `Path`
  - On Windows, starts as a detached process (new process group) and then stores the PID in `stats.json`
  - Fails if the project is already marked running
- Examples:

```bash
dev-env start web
Successfully started 'web'
```

Notes:

- The command string is split by spaces. For complex invocations, consider wrapping in a script/batch file.
- Windows-only detached flags are used via `golang.org/x/sys/windows`.

#### stop

- File: `stop.go`
- Type: `*cobra.Command` (subcommand)
- Use: `dev-env stop <project-name>`
- Behavior:
  - Reads the PID from `stats.json`
  - Windows: uses `taskkill /F /T /PID <pid>`
  - Other OSes: finds the process by PID and calls `Kill()`
  - Removes the project entry from `stats.json`
- Example:

```bash
dev-env stop web
Successfully stopped 'web'
```

---

### Project Structure (files of interest)

- `main.go`: Entrypoint calling `Execute()`
- `root.go`: Defines the root Cobra command
- `init.go`: Implements `init` command
- `list.go`: Implements `list` command and `--name` flag
- `start.go`: Implements `start` command (Windows detached process)
- `stop.go`: Implements `stop` command (taskkill on Windows, Kill() elsewhere)

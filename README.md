## dev-env-cli

A small Go CLI for managing local development projects (start/stop/list) using a simple JSON config in your home directory.

### Features

- Initialize a workspace at $HOME/dev-env-cli
- Track projects in projects.json
- Start commands as detached processes and store their PIDs in stats.json
- Stop running projects by name
- List projects and show RUNNING/STOPPED state
- Lightweight dashboard (web UI) with start/stop controls (embedded via go:embed)

### Requirements

- Go 1.20+
- Windows is the primary target (detached process flags in `start` use Windows APIs). Other platforms may require adjustments.
- Node.js 18+ (only if you want to build the dashboard web UI)

### Install

```bash
go build -o dev-env-cli
# Optionally install to GOPATH/bin
go install github.com/hanzala211/dev-env-cli@latest
```

Add the binary to your PATH if needed.

If you plan to use the dashboard, build the front-end once to generate static assets:

```bash
cd web
npm install
npm run build
# This produces web/dist which will be embedded into the binary and served from memory
```

### Quick Start

```bash
# 1) Initialize workspace under your home directory
dev-env-cli init

# 2) Add projects to $HOME/dev-env-cli/projects.json (see example below; path is optional and the current project directory will be used if not provided)
dev-env-cli add --name --cmd --path

# 3) List projects (the name is optional)
dev-env-cli list --name

# 4) Start a project
dev-env-cli start <project-name>

# 5) Stop a project
dev-env-cli stop <project-name>

# 6) Launch the dashboard (optional web UI)
dev-env-cli dashboard
# Open http://localhost:8080 in your browser
```

### Data files

- `$HOME/dev-env-cli/projects.json`: Array of projects. Expected fields:

  - `name`: Unique project name
  - `path`: Working directory for the command
  - `cmd`: Shell command to start the project (string, split by spaces)

- `$HOME/dev-env-cli/stats.json`: Map of `projectName -> PID` for running projects.

Example `projects.json`:

```json
[
  {
    "name": "web",
    "path": "D:/Work/web-app",
    "cmd": "npm run dev"
  },
  {
    "name": "api",
    "path": "D:/Work/api",
    "cmd": "go run main.go"
  }
]
```

---

### Commands

#### Root command

- File: `root.go`
- Type: `*cobra.Command` (root)
- Description: Sets up the CLI root (`dev-env-cli`) and wires subcommands.

```bash
dev-env-cli --help
```

#### init

- File: `init.go`
- Type: `*cobra.Command` (subcommand)
- Use: `dev-env-cli init`
- Short: Initialize the tool for development
- What it does:
  - Creates `$HOME/dev-env-cli/`
  - Writes empty `projects.json` ([]) and `stats.json` ({}) files
- Example:

```bash
dev-env-cli init
Initialized dev-env-cli in C:/Users/<you>/dev-env-cli
```

#### add

- File: `add.go`
- Type: `*cobra.Command` (subcommand)
- Use: `dev-env-cli add`
- Flags:
  - `--name <project-name>` (required): Project name
  - `--cmd <command>` (required): Command to run (e.g., "npm run dev")
  - `--path <path>` (optional): Project directory (defaults to current directory)
- Behavior:
  - Creates `$HOME/dev-env-cli` if not already initialized
  - Adds a new project entry to `projects.json`
  - You can also provide the command after `--` to include spaces
- Examples:

```bash
dev-env-cli add --name web --cmd "npm run dev" --path D:/Work/web-app
dev-env-cli add --name api -- go run main.go
```

---

#### list

- File: `list.go`
- Type: `*cobra.Command` (subcommand)
- Use: `dev-env-cli list`
- Flags:
  - `--name <project-name>`: Show details for a single project (with path and state)
- Behavior:
  - Reads `projects.json` and `stats.json`
  - Prints `[RUNNING]` or `[STOPPED]` per project
- Examples:

```bash
# List all
dev-env-cli list
[RUNNING] - web
[STOPPED] - api

# Show one by name
dev-env-cli list --name web
[RUNNING] - web - D:/Work/web-app
```

#### start

- File: `start.go`
- Type: `*cobra.Command` (subcommand)
- Use: `dev-env-cli start <project-name>`
- Behavior:
  - Looks up the project by name in `projects.json`
  - Splits `Cmd` by spaces and starts it in `Path`
  - On Windows, starts as a detached process (new process group) and then stores the PID in `stats.json`
  - Fails if the project is already marked running
- Examples:

```bash
dev-env-cli start web
Successfully started 'web'
```

Notes:

- The command string is split by spaces. For complex invocations, consider wrapping in a script/batch file.
- Windows-only detached flags are used via `golang.org/x/sys/windows`.

#### stop

- File: `stop.go`
- Type: `*cobra.Command` (subcommand)
- Use: `dev-env-cli stop <project-name>`
- Behavior:
  - Reads the PID from `stats.json`
  - Windows: uses `taskkill /F /T /PID <pid>`
  - Other OSes: finds the process by PID and calls `Kill()`
  - Removes the project entry from `stats.json`
- Example:

```bash
dev-env-cli stop web
Successfully stopped 'web'
```

#### dashboard

- File: `dashboard.go`
- Type: `*cobra.Command` (subcommand)
- Use: `dev-env-cli dashboard`
- Behavior:
  - Starts a local HTTP server on `http://localhost:8080`
- Serves the embedded dashboard SPA (via go:embed)
  - Exposes REST API endpoints under `/api` to list/start/stop projects
  - Requires `dev-env-cli` to be on your PATH (server shells out to `dev-env-cli start|stop`)
- Example:

```bash
dev-env-cli dashboard
# Then open http://localhost:8080
```

API endpoints (served by the dashboard server):

- `GET /api/projects` → `{ projects: Array<{ name, cmd, path, running }> }`
- `POST /api/projects/start` → body: `{ "name": "<project>" }`
- `POST /api/projects/stop` → body: `{ "name": "<project>" }`

Examples:

```bash
curl http://localhost:8080/api/projects
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"web"}' http://localhost:8080/api/projects/start
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"web"}' http://localhost:8080/api/projects/stop
```

---

### Project Structure (files of interest)

- `main.go`: Entrypoint calling `Execute()`
- `root.go`: Defines the root Cobra command
- `init.go`: Implements `init` command
- `add.go`: Implements `add` command and defines the `Project` type
- `list.go`: Implements `list` command and `--name` flag
- `start.go`: Implements `start` command (Windows detached process)
- `stop.go`: Implements `stop` command (taskkill on Windows, Kill() elsewhere)
- `dashboard.go`: CLI command that starts the dashboard server on port 8080
- `server/server.go`: HTTP server that serves `web/dist` and `/api` routes

---

### Dashboard UI (front-end)

- Location: `web/`
- Build: `npm install && npm run build` to generate `web/dist`
- SPA assets are embedded via go:embed and served from memory
- The UI lists projects and lets you start/stop them via the API

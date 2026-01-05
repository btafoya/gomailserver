# WebUI Control Script Enhancement

## Issue
The `scripts/gomailserver-control.sh` script did not start the WebUI when running in `--dev` mode. The unified WebUI at `web/unified/` needs to be started alongside the Go backend for development work.

## Solution
Enhanced the control script to automatically manage the WebUI development server (Vite) when running in development mode.

## Changes Made

### 1. Configuration Variables Added
- `WEBUI_PID_FILE`: PID file location for WebUI process tracking
- `WEBUI_LOG_FILE`: Log file location for WebUI output
- `WEBUI_DIR`: Directory path to unified WebUI application

### 2. New Functions Added

#### `is_webui_running()`
Checks if the WebUI development server is running by:
- Reading PID from `$WEBUI_PID_FILE`
- Verifying process is active
- Cleaning up stale PID files

#### `start_webui()`
Starts the WebUI development server with:
- Development mode check (only runs in `--dev` mode)
- `pnpm` installation verification
- Automatic `pnpm install` if `node_modules` missing
- Background process startup with logging
- PID file creation for process management
- Health check after startup

#### `stop_webui()`
Gracefully stops the WebUI development server:
- Checks if process is running
- Sends SIGTERM for graceful shutdown
- 10-second timeout before force kill
- PID file cleanup

### 3. Integration Changes

#### `start_server()`
Now calls `start_webui()` after successfully starting the Go backend when in development mode:
```bash
if [ "$MODE" = "development" ]; then
    start_webui
fi
```

#### `stop_server()`
Now calls `stop_webui()` before stopping the Go backend to ensure clean shutdown:
```bash
# Stop WebUI first if it's running
stop_webui
```

#### `show_status()`
Enhanced to display both server and WebUI status:
- Shows WebUI PID when running
- Displays WebUI URL (http://localhost:5173)
- Shows WebUI log file location
- Indicates WebUI only runs in dev mode

#### `show_usage()`
Updated help documentation to include:
- WebUI information in options description
- WebUI URL in examples
- WebUI log file paths
- Development mode behavior explanation

## Testing Results

### Start Server in Dev Mode
```bash
$ ./scripts/gomailserver-control.sh start --dev
[INFO] Mode: Development
[INFO] Using configuration: /home/btafoya/projects/gomailserver/gomailserver.yaml
[INFO] Starting gomailserver...
[SUCCESS] gomailserver started successfully (PID: 2763802)
[INFO] Logs: /home/btafoya/projects/gomailserver/data/gomailserver.log
[INFO] Starting WebUI development server...
[SUCCESS] WebUI started successfully (PID: 2763959)
[INFO] WebUI logs: /home/btafoya/projects/gomailserver/data/webui.log
[INFO] WebUI should be available at http://localhost:5173
```

### Status Check
```bash
$ ./scripts/gomailserver-control.sh status
[SUCCESS] gomailserver is running (PID: 2763802)

    PID    PPID USER     %CPU %MEM     ELAPSED CMD
2763802       1 btafoya   0.5  0.0       00:11 /home/btafoya/projects/gomailserver/build/gomailserver run --config /home/btafoya/projects/gomailserver/gomailserver.yaml

[INFO] Listening ports:
[WARNING] No listening ports found

[SUCCESS] WebUI is running (PID: 2763959)
[INFO] WebUI URL: http://localhost:5173
[INFO] WebUI logs: /home/btafoya/projects/gomailserver/data/webui.log
```

### Stop Server
```bash
$ ./scripts/gomailserver-control.sh stop
[INFO] Stopping WebUI (PID: 2763959)...
[SUCCESS] WebUI stopped
[INFO] Stopping gomailserver (PID: 2763802)...
[SUCCESS] gomailserver stopped
```

### WebUI Accessibility
The WebUI is successfully accessible at http://localhost:5173/admin/ and serves the Vue 3 unified application.

## Development Workflow

### Starting Development Environment
```bash
./scripts/gomailserver-control.sh start --dev
```
This single command now:
1. Starts the Go backend server (gomailserver)
2. Starts the Vue/Vite development server (WebUI)
3. Tracks both processes with PID files
4. Logs both services to separate log files

### Monitoring Both Services
```bash
./scripts/gomailserver-control.sh status
```

### Viewing WebUI Logs
```bash
tail -f data/webui.log
```

### Stopping Everything
```bash
./scripts/gomailserver-control.sh stop
```
This cleanly shuts down both the WebUI and the backend server.

## Production Behavior

In production mode (without `--dev` flag):
- WebUI is NOT started (production uses embedded static files)
- Only the Go backend server runs
- Status command indicates "WebUI is not running (only runs in --dev mode)"

## Requirements

### For WebUI Development
- Node.js (v22.20.0 or compatible)
- pnpm package manager
- WebUI dependencies (auto-installed if missing)

The script automatically:
- Checks for `pnpm` installation
- Runs `pnpm install` if `node_modules` is missing
- Reports clear error messages if requirements are not met

## File Locations

### PID Files
- **Server**: `data/gomailserver.pid`
- **WebUI**: `data/webui.pid`

### Log Files
- **Server**: `data/gomailserver.log`
- **WebUI**: `data/webui.log`

### WebUI Source
- **Directory**: `web/unified/`
- **Dev Server**: Vite (http://localhost:5173)
- **Production Build**: `web/unified/dist/` (served by Go backend)

## Troubleshooting

### WebUI Fails to Start
Check the WebUI log file:
```bash
cat data/webui.log
```

Common issues:
- pnpm not installed: Install with `npm install -g pnpm`
- Port 5173 in use: Stop other Vite dev servers
- Missing dependencies: Script auto-runs `pnpm install`

### WebUI Shows Stale Content
Restart the WebUI to clear Vite cache:
```bash
./scripts/gomailserver-control.sh restart --dev
```

### Both Services Don't Stop
Check for orphaned processes:
```bash
ps aux | grep -E "(vite|pnpm|gomailserver)" | grep btafoya
```

The script includes graceful shutdown with 10-second timeout for WebUI and 30-second timeout for the server.

## Future Enhancements

Potential improvements:
1. Hot reload configuration for the Go backend
2. Parallel startup for faster dev environment initialization
3. Health check endpoints for both services
4. Unified log aggregation and viewing
5. Development mode environment variable injection

## Conclusion

The control script now provides a seamless development experience by managing both the Go backend and Vue frontend with a single command. This eliminates the need to manually start multiple terminals and remember separate commands for each service.

**Development workflow simplified from:**
```bash
# Terminal 1
cd web/unified
pnpm dev

# Terminal 2
./build/gomailserver run --config gomailserver.yaml
```

**To:**
```bash
./scripts/gomailserver-control.sh start --dev
```

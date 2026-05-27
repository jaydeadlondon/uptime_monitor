# Uptime Monitor

Uptime Monitor is a simple service for checking the availability and response times of your services and websites.

## Key Features

- Monitoring of HTTP/HTTPS endpoints
- Outage alerts (configurable via web interface)
- Check history and basic availability metrics
- Separate codebases: backend in Go, frontend in TypeScript (React + Vite)

## Technologies

- TypeScript (frontend)
- Go (backend)
- CSS

## Installation and setup (locally)

Requirements:
- Go 1.20+ (for the backend)
- Node.js 18+ and npm/yarn (for the frontend)

### 1) Clone the repository:

```bash
git clone https://github.com/jaydeadlondon/uptime_monitor.git
cd uptime_monitor
```

### 2) Run the frontend:

```bash
cd frontend
npm install
npm run dev
# or
yarn
yarn dev
```

By default, the frontend runs in development mode at http://localhost:5173 (Vite).

### 3) Run the backend (example):

```bash
# from the root of the repository
# if the backend is located in a subdirectory, navigate to it, for example: cd server
go run ./...
```

The project includes a .env.example file; copy it to .env and specify the necessary environment variables before running.

## Configuration

- Environment variables and check parameters are configured via configuration files or .env (depending on the implementation).
- Frequently used variables: database address, keys for sending notifications, check interval settings.

## Usage

- Via the web interface: add the URL to monitor, set the interval and the trigger threshold.
- View the history of checks and notifications in the interface.

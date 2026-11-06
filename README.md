# Simple Port Scanner

A fast and simple TCP port scanner written in Go. 

I built this project to practice Go concurrency patterns, channels, and network programming using only the standard library.

## What it does

- Scans a range of TCP ports on a target host.
- Uses a **Worker Pool** (goroutines + channels) to scan ports concurrently.
- Supports **Graceful Shutdown**: you can stop the scan at any time by pressing `Ctrl+C`.
- Has zero external dependencies (pure Go standard library).

## How to build

Make sure you have Go installed:

```bash
go build -o portscanner main.go
```

## How to use

Run the binary with flags:

```bash
# Scan local machine (ports 1 to 1024 with 10 workers)
./portscanner -host 127.0.0.1 -start 1 -end 1024 -workers 10

# Scan with custom timeout
./portscanner -host 127.0.0.1 -start 20 -end 100 -timeout 300ms
```

### Available flags

- `-host`: Target IP or domain (default: `127.0.0.1`)
- `-start`: Starting port (default: `1`)
- `-end`: Ending port (default: `1024`)
- `-workers`: Number of concurrent workers (default: `10`)
- `-timeout`: Connection timeout per port (default: `500ms`)

## What I learned

- How to manage a pool of worker goroutines using `sync.WaitGroup`.
- How to send tasks and collect results through Go channels.
- How to handle OS signals (`SIGINT`, `SIGTERM`) with `context` to exit cleanly without leaking resources.
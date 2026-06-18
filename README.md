# Port Scanner

A fast, concurrent TCP port scanner written in Go. It checks which ports are open on a target host, runs every port check in parallel using goroutines, and identifies the service behind the most common ports.

Built as a learning project to practice Go fundamentals: the `flag` package, the `net` package, pointers, maps, and concurrency with goroutines and `sync.WaitGroup`.

## Features

- **Concurrent scanning** — every port is checked in its own goroutine, so a full range is scanned in seconds instead of minutes.
- **Custom port range** — choose exactly which ports to scan with `-start` and `-end`.
- **Full scan** — `-all` scans every port from 1 to 65535.
- **Service detection** — recognises common services (HTTP, SSH, HTTPS, FTP, DNS, MySQL, and more) from a built-in port map.
- **Scan one port** — with `-scan` you can scan ony one port

## Requirements

- [Go](https://go.dev/dl/) 1.18 or newer.

## Installation

Clone the repository and build the executable:

```bash
git clone https://github.com/Giorgioferri/port-scanner.git
cd port-scanner
go build port_scanner.go
```

This produces a standalone binary (`port_scanner` on Linux/macOS, `port_scanner.exe` on Windows) that runs without needing Go installed.

## Usage

| Flag      | Type    | Default     | Description                          |
|-----------|---------|-------------|--------------------------------------|
| `-host`   | string  | `localhost` | Target host to scan                  |
| `-start`  | int     | `1`         | First port of the range             |
| `-end`    | int     | `1`         | Last port of the range              |
| `-all`    | bool    | `false`     | Scan all ports (1–65535)            |
| `-port`    | int    | `0`     | Scan ony one port           |

### Examples

Scan the well-known ports of a host:

```bash
./port_scanner -host scanme.nmap.org -start 1 -end 1024
```

Scan a single specific port:

```bash
./port_scanner -host scanme.nmap.org -start 80 -end 80
```

Scan every port:

```bash
./port_scanner -host scanme.nmap.org -all
```

You can also run it directly without building:

```bash
go run port_scanner.go -host scanme.nmap.org -start 1 -end 1024
```
Scan one single port

```bash
go run port_scanner.go -host scanme.nmap.org -port 80
```

### Sample output

```
port 22 open whit service SSH
port 80 open whit service HTTP
finish
```

Because the scan runs concurrently, open ports may appear in any order — whichever responds first.

## How it works

For each port in the range, the scanner launches a goroutine that attempts a TCP connection with `net.DialTimeout`. A successful connection means the port is open; a failed one (or a timeout) means it is closed or filtered. A `sync.WaitGroup` keeps the program alive until every goroutine has finished. Open ports are then looked up in a `map[int]string` to attach a human-readable service name.

## Roadmap

- [x] Bounded concurrency (worker pool with channels) for reliable full-range scans without hitting the OS socket limit
- [ ] Banner grabbing to read service versions (e.g. `OpenSSH 8.2`)
- [ ] Save results to a file
- [ ] Larger service map, optionally loaded from the official [IANA registry](https://www.iana.org/assignments/service-names-port-numbers)

## Disclaimer

This tool is intended **only** for scanning systems you own or have explicit permission to test. Unauthorised port scanning may be illegal in your jurisdiction. `scanme.nmap.org` is provided by the Nmap project specifically as a legal target for testing scanners.

## License

Released under the MIT License.

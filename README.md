
# noping

`noping`  is a powerful and flexible network diagnostic tool written in Go. It allows users to send ICMP, TCP, or UDP pings to a specified IP address, making it useful for testing connectivity, diagnosing network issues, and measuring response times.

## Features

- Supports ICMP, TCP, and UDP pinging.
- Defaults to ICMP when no port is specified.
- Customizable port, timeout and amount of pings.
- Minimal output mode for clean results.
- Provides ASN data, and geo-location data.
- Fast and efficient.


## Installation

Ensure you have Go installed, then run:
```bash
go install github.com/Bastih18/NoPing@latest
```
Or clone and build manually:
```bash
git clone https://github.com/Bastih18/NoPing.git
cd NoPing
go build -o noping .
```
For local development please follow the [Local Development Guide](https://github.com/Bastih18/NoPing/blob/main/SETUP_LOCAL_DEV.md)
## Usage

```bash
noping <ip> [OPTIONS]
```

#### Arguments
- `ip` - IP address to ping.

#### Options
| Flag | Description |
| - | - |
| `-h, --help` | Show the help menu. |
| `-p, --port <port>` | Specify port to be pinged (default: ICMP if none is given). |
| `-c, --count <count>` | Number of pings (default: 65535) |
| `-t, --timeout <ms>` | Timeout in milliseconds (default: 1000) |
| `-m, --minimal` | Minimal output mode. |
| `-v, --version` | Print detailed version information.
| `--proto <tcp/udp>` | Protocol to use when a port is specified (default: TCP) |
| `--update …[version]` | Update noping to the specified version (when empty, it updates to the latest version).

## Examples

Ping with ICMP:
```bash
noping 192.168.1.1
```

Ping with TCP on port 80:
```bash
noping 192.168.1.1 -p 80
```

Set timeout and number of pings:
```bash
noping 192.168.1.1 -c 10 -t 500
```

Use UDP on port 53:
```bash
noping 192.168.1.1 --proto udp -p 53
```

## License

Licensed under the MIT License. See [LICENSE](https://github.com/Bastih18/NoPing/blob/main/LICENSE) for details.
## Contributing

Contributions are always welcome! Submit issues or pull requests

## Authors

- [@bastih18](https://github.com/bastih18)


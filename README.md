# Automater

Automater is a CLI tool designed to automate common DevOps tasks. It simplifies repetitive operations and helps streamline your workflow.

**Version:** dev (Note: Version is set during build time)

## Installation / Building

To build the `automater` tool, you can use the standard Go build command or the provided Makefile (if configured):

```bash
# Using go build
go build -o automater main.go

# Or using make (if a build target exists in Makefile)
make build
```

This will create an executable file named `automater` in the current directory.

## Usage

```bash
./automater [command]
```

### Available Commands

- **`deploy`**: Handles deployment tasks. (Further documentation needed based on implementation).
- **`monitor`**: Monitors services, applications, hosts, Kubernetes clusters, or system resources.
- **`help`**: Displays help information for any command.

### Global Flags

- `--version`: Display the version of the tool.
- `--help`: Display help information.

## Monitor Command

The `monitor` command helps you track the status and performance of various targets in real-time.

```bash
./automater monitor --service <type> [options...]
```

#### Monitor Flags:

- `--service`, `-s` (string): Specifies the type of monitoring to perform. Required.
  - Available types: `web`, `process`, `cpu`, `memory`, `ping`, `k8s`, `all`.
  - Default: `all`
- `--interval`, `-i` (int): Sets the check interval in seconds.
  - Default: `30`
- `--log`, `-l` (bool): Enables logging output to the console (stderr).
  - Default: `true`
- `--url` (string): The URL to monitor when using `--service web`.
- `--pid` (int): The Process ID (PID) to monitor when using `--service process`.
- `--process-name` (string): The process name to monitor when using `--service process`. (Note: `--pid` takes precedence if both are provided).
- `--host` (string): The hostname or IP address to check reachability for when using `--service ping`.
- `--kubeconfig` (string): Path to the Kubernetes configuration file. Used with `--service k8s`.
  - Default: `~/.kube/config`

#### Monitor Examples:

- **Monitor a website every 15 seconds:**
  ```bash
  ./automater monitor --service web --url https://example.com --interval 15
  ```
- **Monitor a process by name every 5 seconds:**
  ```bash
  ./automater monitor --service process --process-name myapp --interval 5
  ```
- **Monitor CPU usage every 2 seconds:**
  ```bash
  ./automater monitor --service cpu --interval 2
  ```
- **Monitor Memory usage:**
  ```bash
  ./automater monitor --service memory
  ```
- **Check if a host is reachable:**
  ```bash
  ./automater monitor --service ping --host 192.168.1.1
  ```
- **Check Kubernetes cluster connectivity using a specific kubeconfig:**
  ```bash
  ./automater monitor --service k8s --kubeconfig /path/to/my/kube.config
  ```
- **Monitor all supported types (requires relevant flags like --url, --host, etc. to be set for specific checks):**
  ```bash
  ./automater monitor --service all --url https://example.com --process-name nginx --host 8.8.8.8 --interval 60
  ```

## Deploy Command

(Details for the `deploy` command should be added here once implemented.)

## Contributing

(Optional: Add guidelines for contributing if applicable.)

## License

(Optional: Specify the license for the project.)

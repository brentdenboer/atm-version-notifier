# ATM10 Modpack Version Monitor

This application is designed to monitor the version of the "All The Mods 10" (ATM10) modpack running on your server and send a Discord notification via webhook when the version changes.

## Overview

The application works by:

1. **Mounting the ATM10 data volume:** It mounts the `atm10_data` Docker volume (where your ATM10 server files are located) in read-only mode.
2. **Reading the `bcc-common.toml` file:** It reads the `/data/config/bcc-common.toml` file within the mounted volume to extract the `modpackName` and `modpackVersion`.
3. **Saving version information:** It saves the `modpackName` and `modpackVersion` to a reference file in a dedicated volume.
4. **Checking for changes:** It periodically checks if the version in `bcc-common.toml` has changed compared to the saved version by monitoring for file modifications.
5. **Discord Notification:** If the version has changed, it sends a message to a designated Discord channel via a webhook, notifying you of the update.

## Tools Used

- **Go:** The application is written in Go for its excellent concurrency support and robust standard library
- **Docker:** Used for containerization and easy deployment
- **BurntSushi/toml:** Go library for parsing TOML configuration files
- **Discord Webhooks:** For sending notifications about version changes

## Development Setup

Before you begin, ensure you have the following tools installed:

* **Go:** You need Go installed to build the application. You can download it from [https://go.dev/dl/](https://go.dev/dl/). (It's recommended to use a recent version of Go).
* **Docker:** Docker is required to build and run the container. Install Docker Engine and Docker Compose (optional, but recommended) from [https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/).

### Building the Docker Image

1. **Clone the repository:**

    ```bash
    git clone [repository_url] # Replace with your repository URL
    cd [repository_directory] # Replace with your repository directory name
    ```

2. **Build the Docker image:**

    ```bash
    docker build -t atm10-version-monitor .
    ```

## Usage

To run the ATM10 Modpack Version Monitor container, use the following `docker run` command:

```bash
docker run -d \
    --name atm10-version-monitor \
    -v atm10_data:/data:ro \
    -v atm10-version-monitor-reference-data:/reference-data \
    -e DISCORD_WEBHOOK_URL="YOUR_DISCORD_WEBHOOK_URL" \
    -e FILE_CHECK_INTERVAL_SECONDS="60" \
    -e RECHECK_TIMEOUT_SECONDS="30" \
    -e LOG_LEVEL="INFO" \
    atm10-version-monitor
```

### Volume Mounts

The application requires two volume mounts:

1. **ATM10 Data Volume** (`atm10_data:/data:ro`):
   - Contains the modpack configuration files
   - Mounted as read-only for safety
   - Expected to contain `/data/config/bcc-common.toml`

2. **Reference Data Volume** (`atm10-version-monitor-reference-data:/reference-data`):
   - Stores the version reference file
   - Must be writable by the application
   - Contains `/reference-data/version_reference.json`

### Environment Variables

The application can be configured using the following environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DISCORD_WEBHOOK_URL` | Discord webhook URL for notifications | - | Yes |
| `FILE_CHECK_INTERVAL_SECONDS` | Interval between file checks (in seconds) | `60` | No |
| `RECHECK_TIMEOUT_SECONDS` | Timeout after file modification before re-reading (in seconds) | `30` | No |
| `LOG_LEVEL` | Logging level (DEBUG, INFO, WARN, ERROR) | `INFO` | No |
| `CONFIG_FILE_PATH` | Path to the bcc-common.toml file | `/data/config/bcc-common.toml` | No |
| `REFERENCE_FILE_PATH` | Path to the version reference file | `/reference-data/version_reference.json` | No |

### Logging

The application uses a structured logging system with the following features:

- **Log Levels**: Supports multiple log levels (DEBUG, INFO, WARN, ERROR)
- **Output**: Logs are written to both stdout and a log file (`app.log`)
- **Format**: Log entries include timestamp, log level, and source file information
- **Persistence**: Log files are stored in the container and can be accessed via Docker logs

Example log output:
```
2024/03/21 15:04:05 [INFO] main.go:42: ATM10 Modpack Version Monitor started
2024/03/21 15:04:05 [DEBUG] main.go:50: File check interval: 60 seconds, Recheck timeout: 30 seconds
2024/03/21 15:04:06 [INFO] main.go:65: Version change detected for ATM10! Old: 0.1.0, New: 0.1.1
```

You can view the logs using:
```bash
# View container logs
docker logs atm10-version-monitor

# Follow log output
docker logs -f atm10-version-monitor
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

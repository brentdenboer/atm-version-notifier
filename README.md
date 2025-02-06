# ATM10 Modpack Version Monitor

This application monitors the "All The Mods 10" (ATM10) modpack version on a server and sends a Discord notification via webhook when the version changes. It's available on Docker Hub as `brentdboer/atm10-version-notifier`.

## Overview

The application's core functionality is as follows:

1.  **Configuration:** Loads settings (Discord webhook URL, file paths) from environment variables.
2.  **Version Extraction:** Reads the `modpackName` and `modpackVersion` from `/data/config/bcc-common.toml` within a mounted volume.
3.  **Persistent Reference:** Stores the last known version in `/reference-data/version_reference.json`.
4.  **Change Detection:** Monitors `bcc-common.toml` for modifications.  On change, it re-reads the file and compares the new version to the reference.
5.  **Discord Notifications:** Sends a Discord message via webhook if a version difference is detected.
6. **Continuous Monitoring**: using `fsnotify`.
7. **Robust Error Handling**: Descriptive error messages and context wrapping (`%w`).
8. **Structured Logging**:  Uber Zap logging to stdout (console-friendly) and a log file (`app.log`, JSON format).  Sensitive data is masked.
9. **Atomic Reference Updates**: Preventing data corruption.
10. **Input Validation**:  Validates environment variables, TOML configuration, and reference file data.
11. **Debouncing**: Prevents multiple notifications for rapid file changes.
12. **Secure Dockerfile**: Minimal image, non-root user, `ca-certificates`.

## Tools Used

*   **Go:** Programming language.
*   **Docker:** Containerization.
*   **BurntSushi/toml:** TOML parsing.
*   **fsnotify/fsnotify:** File system notifications.
*   **Uber Zap:** Structured logging.
*   **Discord Webhooks:** For notifications.

## Deployment (Docker Hub)

The easiest way to deploy the application is to use the pre-built image from Docker Hub: `brentdboer/atm10-version-notifier`.

### Docker Run Command

```bash
docker run -d \
    --name atm10-version-monitor \
    -v atm10_data:/data:ro \
    -v atm10-version-monitor-reference-data:/reference-data \
    -e DISCORD_WEBHOOK_URL="YOUR_DISCORD_WEBHOOK_URL" \
    -e LOG_LEVEL="INFO" \
    brentdboer/atm10-version-notifier
```

**Explanation:**

*   `-d`: Runs in detached mode.
*   `--name`: Assigns a name.
*   `-v atm10_data:/data:ro`: Mounts your ATM10 data volume (read-only).  **Important:** This volume *must* contain `/data/config/bcc-common.toml`.
*   `-v atm10-version-monitor-reference-data:/reference-data`: Mounts a volume for the persistent `version_reference.json` file (read-write).
*   `-e DISCORD_WEBHOOK_URL="..."`: Sets the required Discord webhook URL. *Replace `YOUR_DISCORD_WEBHOOK_URL` with your actual URL.*
*   `-e LOG_LEVEL="INFO"`: Sets the logging level (INFO, DEBUG, WARN, ERROR).
*   `brentdboer/atm10-version-notifier`:  Pulls and runs the image from Docker Hub.

### Volume Mounts

*   **`atm10_data:/data:ro` (Read-Only):** Your ATM10 server data.  The app expects `/data/config/bcc-common.toml` here.
*   **`atm10-version-monitor-reference-data:/reference-data` (Read-Write):** Stores the `version_reference.json` file (last known version).

### Environment Variables

| Variable              | Description                                        | Default                                  | Required |
| --------------------- | -------------------------------------------------- | ---------------------------------------- | -------- |
| `DISCORD_WEBHOOK_URL` | Your Discord webhook URL.                          | *None*                                   | Yes      |
| `LOG_LEVEL`           | Logging level (DEBUG, INFO, WARN, ERROR, FATAL).   | `INFO`                                   | No       |
| `CONFIG_FILE_PATH`    | Path to `bcc-common.toml` *within the container*.  | `/data/config/bcc-common.toml`           | No       |
| `REFERENCE_FILE_PATH` | Path to the reference file *within the container*. | `/reference-data/version_reference.json` | No       |

### Logging

*   **Levels:** DEBUG, INFO, WARN, ERROR, FATAL (controlled by `LOG_LEVEL`).
*   **Output:**
    *   Standard Output (stdout): Human-readable.
    *   `app.log` (in the container): JSON format.
*   **Structure:** Includes fields like `time`, `level`, `caller`, `msg`.
*   **Security:** Masks sensitive information (webhook URLs).

**Viewing Logs:**

```bash
docker logs -f atm10-version-monitor  # Real-time
docker logs atm10-version-monitor     # Past logs
```

## Development Setup (Optional)

If you want to build the image yourself or contribute to the project:

**Prerequisites:**

*   Go (1.21+): [https://go.dev/dl/](https://go.dev/dl/)
*   Docker: [https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/)

**Building:**

1.  Clone the repository:
    ```bash
    git clone [repository_url]
    cd [repository_directory]
    ```
2.  Build the Docker image:
    ```bash
    docker build -t atm10-version-monitor .
    ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file.

## Contributing

Contributions are welcome!  Submit a Pull Request with well-documented, error-handled Go code.

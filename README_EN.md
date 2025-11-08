# gs-write

English | [日本語](README.md)

A simple CLI tool that writes standard input to a new Google Spreadsheet.
Works seamlessly with pipes (`|`) to connect with other commands.

![Go](https://img.shields.io/badge/Go-1.24-blue.svg)

## Overview

`gs-write` is a command-line tool that outputs the results of commands like `ls -l`, `cat report.csv`, or `ps aux` directly to a new Google Spreadsheet.
A new spreadsheet is created with each execution, and the URL of the created sheet is returned to standard output upon completion. This makes it easy to integrate data in shell scripts or copy the resulting URL to the clipboard.

## Features

- Reads CSV data from standard input and writes to Google Spreadsheet
- Always creates a new spreadsheet (no append functionality to existing sheets)
- Seamless integration with other commands via pipes (`|`)
- Freely specify spreadsheet title with `--title` option
  - If no title is specified, it's automatically generated from the execution timestamp (`YYYYMMDDHHMMSS+gs`)
- Freeze rows and columns with `--freeze-rows` and `--freeze-cols` options
- Set basic filter with `--filter-header-row` option
- Returns the URL of the created spreadsheet to standard output on success

## Installation

### Download Pre-built Binary (Recommended)

Download the appropriate file for your platform from the [Releases page](https://github.com/foresukecom/gs-write/releases).

#### macOS

**For Intel Mac:**
```bash
# Download
curl -LO https://github.com/foresukecom/gs-write/releases/latest/download/gs-write_Darwin_x86_64.tar.gz
# Extract
tar xzf gs-write_Darwin_x86_64.tar.gz
# Move to executable location
sudo mv gs-write /usr/local/bin/
```

**For Apple Silicon (M1/M2/M3):**
```bash
# Download
curl -LO https://github.com/foresukecom/gs-write/releases/latest/download/gs-write_Darwin_arm64.tar.gz
# Extract
tar xzf gs-write_Darwin_arm64.tar.gz
# Move to executable location
sudo mv gs-write /usr/local/bin/
```

#### Linux

```bash
# Download
curl -LO https://github.com/foresukecom/gs-write/releases/latest/download/gs-write_Linux_x86_64.tar.gz
# Extract
tar xzf gs-write_Linux_x86_64.tar.gz
# Move to executable location
sudo mv gs-write /usr/local/bin/
```

#### Windows

1. Download [gs-write_Windows_x86_64.zip](https://github.com/foresukecom/gs-write/releases/latest/download/gs-write_Windows_x86_64.zip)
2. Extract the ZIP file
3. Place `gs-write.exe` in your desired folder
4. Add to PATH if needed

### Install via Go

If you have a Go environment set up, you can install with the following command:

```bash
go install github.com/foresukecom/gs-write@latest
```

### Build from Source

Clone this repository and build:

```bash
git clone https://github.com/foresukecom/gs-write.git
cd gs-write
go build -o gs-write .
```

## Setup

### 1. Create a project in Google Cloud Console

1. Access [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project
3. Enable Google Sheets API
4. Create OAuth 2.0 Client ID (Application type: Desktop app)
5. Download `credentials.json`

### 2. Authentication

When using for the first time, or when authentication is required, run the following command:

#### Method 1: Interactively enter credentials

```bash
gs-write auth
```

Follow the prompts to paste the contents of `credentials.json`.

#### Method 2: Load credentials from file

```bash
gs-write auth --credentials ./credentials.json
```

#### Authentication Flow

When you run the command, follow these steps to authenticate:

1. **An authorization URL will be displayed**
   ```
   Please visit the following URL to authorize this application:
   https://accounts.google.com/o/oauth2/auth?...
   ```

2. **Open the URL in your browser**
   - Open the displayed URL in your browser
   - Sign in with your Google account

3. **Authorize the application**
   - You'll see a screen saying "[App name] wants to access your Google Account"
   - Click "Continue"

4. **Get the authorization code**
   - After authorization, you'll be redirected to `http://localhost`, which typically shows an `ERR_CONNECTION_REFUSED` error
   - **This is expected behavior**
   - Copy the code from the URL in your browser's address bar after `code=`
   - Example: For `http://localhost/?code=4/0AbCD...XYZ&scope=...`, copy `4/0AbCD...XYZ`

5. **Enter the authorization code**
   ```
   Enter the authorization code:
   ```
   - Return to your terminal and paste the copied authorization code, then press Enter

6. **Authentication complete**
   ```
   Authentication successful!
   Authentication saved to: ~/.config/gs-write/auth.json
   ```

When authentication is successful, credentials and tokens are saved to `~/.config/gs-write/auth.json`.

## Usage

### Basic Usage

```bash
# Convert CSV file to spreadsheet
cat data.csv | gs-write

# Save command output to spreadsheet
ls -l | gs-write

# Create spreadsheet with specified title
ps aux | gs-write --title "Process List $(date +%Y%m%d)"

# Copy URL to clipboard (macOS)
cat report.csv | gs-write | pbcopy

# Copy URL to clipboard (Linux with xclip)
cat report.csv | gs-write | xclip -selection clipboard
```

### Freeze Panes

You can freeze header rows or specific columns:

```bash
# Freeze the first row (header row)
cat data.csv | gs-write --freeze-rows 1

# Freeze the first 2 columns
cat data.csv | gs-write --freeze-cols 2

# Freeze both the first row and first column
cat data.csv | gs-write --freeze-rows 1 --freeze-cols 1

# Combine title and freeze panes
cat employee.csv | gs-write --title "Employee List" --freeze-rows 1
```

### Basic Filter

You can set a basic filter on the data:

```bash
# Set filter with row 1 as header
cat data.csv | gs-write --filter-header-row 1

# Set filter with row 2 as header
cat data.csv | gs-write --filter-header-row 2

# Combine freeze panes and filter
cat data.csv | gs-write --freeze-rows 1 --filter-header-row 1

# Combine all options
cat employee.csv | gs-write --title "Employee List" --freeze-rows 1 --filter-header-row 1
```

### Character Encoding

You can handle CSV files with encodings other than UTF-8, such as Shift_JIS (SJIS) or EUC-JP:

```bash
# Read CSV file with Shift_JIS (SJIS) encoding
cat data_sjis.csv | gs-write --encoding sjis

# Read CSV file with EUC-JP encoding
cat data_eucjp.csv | gs-write --encoding euc-jp

# UTF-8 (default) doesn't require specification
cat data_utf8.csv | gs-write

# Combine encoding with other options
cat data_sjis.csv | gs-write --encoding sjis --title "Employee List" --freeze-rows 1
```

Supported encodings:
- `utf-8` (default): UTF-8 encoding
- `sjis`: Shift_JIS (standard Japanese encoding on Windows)
- `euc-jp`: EUC-JP (Japanese encoding on Unix-like systems)

### Options

- `--title <title>`: Specify the spreadsheet title. If not specified, it's automatically generated from the timestamp.
- `--freeze-rows <number>`: Freeze the specified number of rows from the top. Overrides config file value.
- `--freeze-cols <number>`: Freeze the specified number of columns from the left. Overrides config file value.
- `--filter-header-row <row-number>`: Set basic filter with the specified row as header. Overrides config file value.
- `--encoding <encoding>`: Specify the character encoding of input CSV (`utf-8`, `sjis`, `euc-jp`). Default is `utf-8`.

### Configuration File

Frequently used configuration values (freeze rows/columns, filter header row, etc.) can be saved in a configuration file.

The configuration file is saved at `~/.config/gs-write/config.toml` and can be managed with the `gs-write config` command.

When command-line options are specified, they override the configuration file values (Priority: CLI > Config file > Default).

```bash
# Display current configuration
gs-write config list

# Get specific configuration value
gs-write config get freeze.rows
gs-write config get filter.header_row

# Change configuration value
gs-write config set freeze.rows 1
gs-write config set freeze.cols 2
gs-write config set filter.header_row 1

# Delete configuration value (revert to default)
gs-write config unset freeze.rows
gs-write config unset filter.header_row
```

#### Available Configuration Settings

- `freeze.rows`: Number of rows to freeze (default: 0)
- `freeze.cols`: Number of columns to freeze (default: 0)
- `filter.header_row`: Filter header row number (default: 0 = no filter)

### Subcommands

#### `gs-write auth`

Authenticate with Google Sheets API.

```bash
# Authenticate interactively
gs-write auth

# Load credentials from file
gs-write auth --credentials ./credentials.json
```

#### `gs-write config`

Manage configuration file.

```bash
# Display configuration list
gs-write config list

# Get specific configuration value
gs-write config get freeze.rows

# Change configuration value
gs-write config set freeze.rows 1

# Delete configuration value
gs-write config unset freeze.rows
```

#### `gs-write version`

Display version information.

```bash
gs-write version
```

## Data Format

`gs-write` expects CSV format data. Data read from standard input is parsed as comma-separated (`,`).

### Example

```bash
echo "Name,Age,City
Alice,30,Tokyo
Bob,25,Osaka" | gs-write --title "User List"
```

## Project Structure

```
.
├── README.md           # This file (Japanese)
├── README_EN.md        # This file (English)
├── cmd/                # Cobra command definitions
│   ├── auth.go         # Auth command
│   ├── config.go       # Config command
│   ├── root.go         # Root command (main functionality)
│   └── version.go      # Version command
├── pkg/                # Internal packages
│   ├── auth/           # Authentication logic
│   │   └── auth.go
│   ├── config/         # Configuration management
│   │   └── config.go
│   └── sheets/         # Google Sheets API client
│       └── sheets.go
├── go.mod              # Go Modules
├── go.sum              # Go Modules checksum
└── main.go             # Entry point
```

## Configuration File Locations

gs-write uses the following files:

- `~/.config/gs-write/auth.json` - OAuth 2.0 credentials and token
- `~/.config/gs-write/config.toml` - User settings (freeze.rows, freeze.cols, etc.)

## Troubleshooting

### Authentication Error

The authentication token may have expired. Re-authenticate:

```bash
gs-write auth
```

### API Quota Error

Google Sheets API has usage limits. If you're sending a large number of requests, wait a moment and try again.

## License

MIT License

## Contributing

Issues and Pull Requests are welcome!

## Development

### Development Environment

This project can be developed using VS Code Dev Containers.

1. Open the project in VS Code
2. Select "Reopen in Container"
3. Start development in the container

### Build

```bash
go build -o gs-write .
```

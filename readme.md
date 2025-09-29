# AirType CLI

A command-line tool that allows you to type on your iPhone from your PC using the [AirType](https://apps.apple.com/us/app/airtype-type-from-your-computer/id922932291) application.

This rewritten version is built in Go and offers a more robust, cross-platform experience with automatic reconnection capabilities.

## Features

- **Interactive Typing**: Run the tool in your terminal and type directly to your device.
- **File-Based Typing**: Automatically type the content of a text file.
- **Cross-Platform**: Works on Windows, macOS, and Linux.
- **Automatic Reconnection**: If the connection to your device drops, the tool will automatically try to reconnect.

## How It Works

The AirType application creates a small web server on your iPhone. When you select AirType as your keyboard, it directs you to a local URL where you can type in a text field on the page. This tool connects to the underlying WebSocket service to send keystrokes from your computer to your iPhone.

## Project Structure

The project is organized into two main directories:
- `cmd/`: Contains the entry points for the two command-line tools (`airtype` and `typetext`).
- `internal/`: Holds the shared `airtype` library, which manages the WebSocket connection and communication logic.

## IOs installation
To use this tool, you need to have the AirType app installed on your iPhone. You can download it from the [App Store](https://apps.apple.com/us/app/airtype-type-from-your-computer/id922932291).
This is a free app that works as a keyboard extension.
Now open settings, go to general, then keyboard, and finally keyboards. Add a new keyboard and under third party keyboards, select AirType. After adding it, tap on it and enable "Allow Full Access".

## PC Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/denizsincar29/airtype.git
   ```

2. **Navigate to the project directory:**
   ```bash
   cd airtype
   ```

3. **Install dependencies:**
   ```bash
   go mod tidy
   ```

## Usage

Before running either tool, you will be prompted to enter the IP address displayed on the AirType keyboard extension on your iPhone. The IP will be saved to an `ip.txt` file for future use.

**Note:** Ensure your iPhone and PC are on the same network.

### Interactive Typing

To type interactively from your terminal:
```bash
go run ./cmd/airtype
```
Press `Esc` or `Ctrl+C` to exit.

### Typing from a File

To automatically type the contents of a file:
1. Create a text file (e.g., `text.txt`) with the content you want to send.
2. Run the `typetext` command:
   ```bash
   go run ./cmd/typetext --file text.txt
   ```
   You can use the `--file` flag to specify a different input file.

## Disclaimer

This is not an official application. It was created for fun and learning purposes. Use it at your own risk.

Since it uses an open and straightforward method, it is easy to read the JavaScript from the page and create your own client. Therefore, do not expect any security or privacy from this project. Use it only on trusted networks, as anyone on the same network can easily sniff the traffic and see what you type.

## License

This project is licensed under the MIT License.
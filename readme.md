# AirType CLI

A cross-platform command-line tool that lets you type on your iPhone from your PC using the [AirType](https://apps.apple.com/us/app/airtype-type-from-your-computer/id922932291) app.

Built in **Go** for robustness, automatic reconnection, and a smooth typing experience.

---

## Quick Start

```bash
git clone https://github.com/denizsincar29/airtype.git
cd airtype
go build -o airtype ./cmd/airtype
./airtype --ip 1.25   # or just ./airtype if ip.txt exists
```

---

## Features

* **Interactive Typing** – Type directly from your terminal.
* **File-Based Typing** – Send file contents with `--file`.
* **Clipboard Support** – Send clipboard text with `-c`.
* **Cross-Platform** – Works on Windows, macOS, and Linux.
* **Automatic Reconnection** – Reconnects if connection drops.

---

## How It Works

The AirType app starts a local web server on your iPhone.
When you switch to the AirType keyboard, it shows an IP address.
This CLI connects to its **WebSocket** endpoint and streams your keystrokes.

---

## Project Structure

* `cmd/` – CLI entry point (`airtype`).
* `airtype/` – Library that manages the WebSocket connection and communication.

---

## iOS Setup

1. Install [AirType](https://apps.apple.com/us/app/airtype-type-from-your-computer/id922932291) (free) from the App Store.
2. Open **Settings → General → Keyboard → Keyboards**.
3. Tap **Add New Keyboard**, scroll down to **Third-Party Keyboards**, and select **AirType**.
4. Tap on **AirType** in the list and enable **Allow Full Access**.

---

## PC Installation

1. Install [Go](https://go.dev/dl/).
2. Clone the repository:

   ```bash
   git clone https://github.com/denizsincar29/airtype.git
   cd airtype
   ```
3. Install dependencies:

   ```bash
   go mod tidy
   ```
4. (Optional) Build a binary:

   ```bash
   go build -o airtype ./cmd/airtype
   ```

---

## Usage

Switch to the **AirType keyboard** on your iPhone.
The keyboard will display an IP address — use it with `--ip`:

```bash
./airtype --ip <ADDRESS>
```

If `--ip` is omitted, the tool will read from `ip.txt` in the current directory automatically:

```bash
./airtype    # uses ip.txt by default
```

### `--ip` Flag

Supports:

* Full IP (`192.168.1.25`)
* Last two octets (`1.25` → expands to `192.168.1.25`)
* Filename containing IP (`ip.txt`)

---

### Interactive Typing

```bash
./airtype --ip <ADDRESS>
```

Press **Ctrl+C** (or Esc in some terminals) to exit.

---

### Typing from a File

```bash
./airtype --file mytext.txt --ip <ADDRESS>
```

or, if `ip.txt` exists:

```bash
./airtype --file mytext.txt
```

---

### Typing from Clipboard

```bash
./airtype -c --ip <ADDRESS>
```

or:

```bash
./airtype -c
```

---

## Notes

* Your iPhone and PC must be on the same Wi-Fi network.
* Typing speed may vary depending on network latency.

---

## Disclaimer

This is an **unofficial** project made for fun and learning.
Traffic is sent over an **unencrypted WebSocket**, so use only on trusted networks.

---

## License

Licensed under the MIT License.

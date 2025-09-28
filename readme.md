# AirType CLI

A command-line tool that allows you to type on your iPhone from your PC using the [AirType](https://apps.apple.com/us/app/airtype-type-from-your-computer/id922932291) application.

## How It Works

The AirType application creates a small web server on your iPhone. When you select AirType as your keyboard, it directs you to a local URL where you can type in a text field on the page.

This project analyzes the page's source code and uses simple WebSockets to send the text you type on your computer to your iPhone. This command-line tool simplifies the process.

The tool was developed using various programming languages I have experience with. I started with Python, then moved to Rust, and finally settled on Go.

## Installation and Usage

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

4. **Run the application:**
   ```bash
   go run .
   ```

You will then be prompted to enter the IP address (without the port) displayed on the AirType extension. Once you enter it, you can start typing on your iPhone from your PC.

**Note:** Ensure your iPhone and PC are on the same network.

## Disclaimer

This is not an official application. It was created for fun and learning purposes. Use it at your own risk.

Since it uses an open and straightforward method, it is easy to read the JavaScript from the page and create your own client. Therefore, do not expect any security or privacy from this project. Use it only on trusted networks, as anyone on the same network can easily sniff the traffic and see what you type.

## License

This project is licensed under the MIT License.
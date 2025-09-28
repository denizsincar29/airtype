# Airtype
A command line tool that allows you to type on your iPhone from your PC using [AirType](https://apps.apple.com/us/app/airtype-type-from-your-computer/id922932291) application.

## How it works
The AirType application creates a tiny web server on your iPhone. When you choose AirType as your keyboard, it asks you to browse to it's local URL and start typing inside the text field right on the page.
I've analyzed the page source and found out that it uses very simple websockets to send the text you type on your computer to the iPhone. So I decided to make it all easier by moving the process into a command line tool.
The tool was made using every languages i programmed in :) actually I started with Python, than rust, than go and stopped at it.

## Installation and usage
```bash
git clone git@github.com:denizsincar29/airtype.git
cd airtype
go mod tidy
go run .
```
Than you will be asked to enter the IP address (without the port) that is displayed on the AirType extension. After you enter it, you will be able to type on your iPhone from your PC.
Ah, I forgot! Make sure your iPhone and PC are on the same network.

## Note
This is not an official application. I made it just for fun and learning purposes. Use it at your own risk.
Since it uses an open and welcoming method, it is a piece of cake to just read the js from the page and make your own client. So don't expect any security or privacy from this project.
Use it only in trusted networks, because anyone in the same network can easily sniff the traffic and see what you type.

## License
This project is licensed under the MIT License.
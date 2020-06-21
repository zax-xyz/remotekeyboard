# KDE Connect Remote Keyboard

This is a simple way to interact with the KDE Connect "remote keyboard" plugin. It displays a small window into which inputs are made, and mirrored to the remote device. No text field is shown, this may be added as an option in the future. Personally, I didn't feel it was necessary as I can simply look at the other device's display.

## Usage

```sh
./remotekeyboard DEVICE_ID
```

The first argument given to the program is used as the device ID (used by KDE Connect) to send the inputs to. If no argument is given, it tries to find the first available device that has the remote keyboard plugin enabled.

If no devices have the plugin enabled, or the device passed in arguments doesn't have the plugin enabled, then the program will quit with the error message `Remote keyboard plugin not enabled!`.

An interface with which to choose a device from a list view may be added in the future.

## Installing

A pre-built binary can be found from the [releases](https://github.com/zaxutic/kdeconnect-remotekeyboard/releases) page, and added to your `PATH`. Alternatively, you can build from source.

### Building From Source

This requires `go` to be installed.

```sh
# Get dependencies
go get ./...

# Build with reduced file size
make
sudo make install

# or without
go build remotekeyboard.go
sudo make install
```

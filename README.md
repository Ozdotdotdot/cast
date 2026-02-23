# cast

A CLI tool for sending text-to-speech messages to Google Home and Chromecast devices on your local network.

```
cast "Dinner is ready"
cast bedroom "Goodnight"
cast bedroom,kitchen "Meeting in 5 minutes"
```

## Install

```bash
go install github.com/Ozdotdotdot/cast@latest
```

Make sure `$GOPATH/bin` (default `~/go/bin`) is in your `PATH`.

## Setup

On first run, `cast` will scan your network for devices automatically. You can re-scan anytime:

```bash
cast discover
```

Devices are saved to `~/.config/cast/config.toml`.

## Usage

```bash
# Speak to the default device
cast "Hello world"

# Speak to a specific device
cast kitchen "Timer is done"

# Speak to multiple devices
cast kitchen,bedroom "Good morning"

# Speak to a device group
cast everywhere "Attention please"
```

## Commands

| Command          | Description                          |
|------------------|--------------------------------------|
| `cast <message>` | Speak to the default device          |
| `cast discover`  | Scan for devices on the network      |
| `cast devices`   | List configured devices and aliases  |
| `cast groups`    | List configured device groups        |
| `cast config`    | Show config file location            |

## How it works

1. Text is converted to speech via Google Translate TTS
2. A temporary local HTTP server serves the audio file
3. The Chromecast protocol loads and plays the audio on the target device(s)

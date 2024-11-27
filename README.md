# rmEventTrigger

Do some actions based on events from the reMarkable tablet. This application listens in to the output of `journalctl` and based on that takes some actions.

Current implemented events:
* Sync completion (only happens if WiFi is connected)

## Installation

You can execute this on the reMarkable, to fetch the application and install it as a service.

```bash
wget -O - https://raw.githubusercontent.com/crazyfacka/rmEventTrigger/refs/heads/main/install.sh | bash
```

## Configuration

You need to have a `JSON` file, in the following format:

```json
{"conf": [
    {
        "Event": "Sync",
        "Actions": ["/bin/ls -la"]
    },
    ...
]}
```

If you have multiple entries with the same Event, the application will concatenate all actions and execute them sequentially.

## Build

Requirements:
* Go >= 1.23

How:
```bash
$ make
```

## Run

Although it is preferable to have this running as a service, you can run it standalone to test or validate any points.

```bash
$ ./app -c conf.json
```

## Notes

Currently the `Makefile`, and the binary on the Releases page, it's for the reMarkable Paper Pro, which runs on ARM64 architecture. It should be easy to adapt to the reMarkable 2 (and perhaps the 1).
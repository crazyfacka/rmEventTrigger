#!/bin/sh

BINARY_URL="https://github.com/crazyfacka/rmEventTrigger/releases/download/v0.1/app"
SERVICE_FILE_URL="https://raw.githubusercontent.com/crazyfacka/rmEventTrigger/refs/tags/v0.1/rm-event-trigger.service"
BINARY_NAME="app"
SERVICE_NAME="rm-event-trigger"

BIN_DIR="/home/root/rmeventtrigger"
SERVICE_DIR="/etc/systemd/system"

echo "Installing rmEventTrigger"

if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

mkdir -p "$BIN_DIR"

echo "Downloading binary..."
wget -O "$BIN_DIR/$BINARY_NAME" "$BINARY_URL"
if [ $? -ne 0 ]; then
    echo "Failed to download binary."
    exit 1
fi

chmod +x "$BIN_DIR/$BINARY_NAME"
if [ $? -ne 0 ]; then
    echo "Failed to make binary executable."
    exit 1
fi

echo "Downloading service file..."
wget -O "$SERVICE_DIR/$SERVICE_NAME.service" "$SERVICE_FILE_URL"
if [ $? -ne 0 ]; then
    echo "Failed to download service file."
    rm -f "$BIN_DIR/$BINARY_NAME"
    exit 1
fi

systemctl daemon-reload
if [ $? -ne 0 ]; then
    echo "Failed to reload systemd daemon."
    rm -f "$BIN_DIR/$BINARY_NAME"
    rm -f "$SERVICE_DIR/$SERVICE_NAME.service"
    exit 1
fi

echo "Enabling and starting the service..."
systemctl enable "$SERVICE_NAME.service"
if [ $? -ne 0 ]; then
    echo "Failed to enable service."
    rm -f "$BIN_DIR/$BINARY_NAME"
    rm -f "$SERVICE_DIR/$SERVICE_NAME.service"
    exit 1
fi

systemctl start "$SERVICE_NAME.service"
if [ $? -ne 0 ]; then
    echo "Failed to start service."
    systemctl disable "$SERVICE_NAME.service"
    rm -f "$BIN_DIR/$BINARY_NAME"
    rm -f "$SERVICE_DIR/$SERVICE_NAME.service"
    exit 1
fi

echo "Installation and service setup completed successfully!"
exit 0

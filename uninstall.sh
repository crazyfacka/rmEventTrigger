#!/bin/sh

SERVICE_NAME="rm-event-trigger"
SERVICE_DIR="/etc/systemd/system"

echo "Removing rmEventTrigger service."
echo "Binaries and app folder will be left untouched."

if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

echo "Stopping and disabling the service..."
systemctl stop "$SERVICE_NAME.service"
if [ $? -ne 0 ]; then
    echo "Failed to stop service."
    exit 1
fi

systemctl disable "$SERVICE_NAME.service"
if [ $? -ne 0 ]; then
    echo "Failed to disable service."
    exit 1
fi

echo "Removing service file..."
rm -f "$SERVICE_DIR/$SERVICE_NAME.service"
if [ $? -ne 0 ]; then
    echo "Failed to remove service file."
    exit 1
fi

systemctl daemon-reload
if [ $? -ne 0 ]; then
    echo "Failed to reload systemd daemon."
    exit 1
fi

echo "Uninstallation completed successfully!"
exit 0
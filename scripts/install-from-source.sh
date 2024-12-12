#!/bin/bash

echo "Building xdir from source..."
go build -o xdir

echo "Installing xdir to /usr/local/bin..."
sudo mv xdir /usr/local/bin/

echo "Installation complete. You can now use 'xdir' command."

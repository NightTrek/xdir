#!/bin/bash

# Create releases directory if it doesn't exist
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
RELEASES_DIR="$SCRIPT_DIR/releases"
mkdir -p "$RELEASES_DIR"

# Versions of Go to target
VERSION="1.21"

# Platforms to target
PLATFORMS=("windows/amd64" "windows/386" "darwin/amd64" "darwin/arm64" "linux/amd64" "linux/386" "linux/arm" "linux/arm64")

# Build for each platform
for PLATFORM in "${PLATFORMS[@]}"
do
    # Split platform into OS and architecture
    IFS='/' read -r -a array <<< "$PLATFORM"
    GOOS="${array[0]}"
    GOARCH="${array[1]}"
    
    # Set output binary name based on platform
    if [ $GOOS = "windows" ]; then
        OUTPUT_NAME="xdir.exe"
    else
        OUTPUT_NAME="xdir"
    fi
    
    # Create platform-specific directory
    OUTPUT_DIR="$RELEASES_DIR/${GOOS}-${GOARCH}"
    mkdir -p "$OUTPUT_DIR"
    
    echo "Building for $GOOS/$GOARCH..."
    
    # Build binary
    cd "$SCRIPT_DIR/.." && GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUTPUT_DIR/$OUTPUT_NAME"
    
    # Create archive
    if [ $GOOS = "windows" ]; then
        zip -j "$OUTPUT_DIR/xdir-${GOOS}-${GOARCH}.zip" "$OUTPUT_DIR/$OUTPUT_NAME"
    else
        tar -czf "$OUTPUT_DIR/xdir-${GOOS}-${GOARCH}.tar.gz" -C "$OUTPUT_DIR" "$OUTPUT_NAME"
    fi
done

echo "Build complete! Release archives are in: $RELEASES_DIR"
ls -la "$RELEASES_DIR"/*/*

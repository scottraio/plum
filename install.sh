#!/bin/bash

# This script installs the Plum framework and renames the binary to "plum".

# Display welcome message and ASCII art


# Install the Plum framework
echo "Installing the Plum..."
go install github.com/scottraio/plum-scaffold@latest

# Rename the binary to "plum"
echo "Renaming the binary to 'plum'..."
mv $GOPATH/bin/plum-scaffold $GOPATH/bin/plum

echo "plum has been successfully installed!"

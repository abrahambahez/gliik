#!/bin/bash

VERSION=$(git describe --tags --always --dirty)
go build -ldflags "-X 'github.com/yourusername/gliik/cmd.version=$VERSION'" -o gliik

echo "Built gliik with version: $VERSION"

cp gliik "$HOME/.local/bin/gliik"
echo "Installed to $HOME/.local/bin/gliik"

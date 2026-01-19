#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

echo -e "${BLUE}╔═══════════════════════════════╗${NC}"
echo -e "${BLUE}║   bubblefetch Installer      ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════╝${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go is not installed${NC}"
    echo -e "  Please install Go from https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✓ Go is installed${NC}"

# Build the binary
echo -e "${YELLOW}Building bubblefetch...${NC}"
go build -ldflags="-s -w" -o bubblefetch ./cmd/bubblefetch

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build successful${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

# Install to system
INSTALL_DIR="/usr/local/bin"
if [ -w "$INSTALL_DIR" ]; then
    echo -e "${YELLOW}Installing to $INSTALL_DIR...${NC}"
    mv bubblefetch "$INSTALL_DIR/"
    ln -sf "$INSTALL_DIR/bubblefetch" "$INSTALL_DIR/bf"
    echo -e "${GREEN}✓ Installed successfully${NC}"
else
    echo -e "${YELLOW}Installing to $INSTALL_DIR (requires sudo)...${NC}"
    sudo mv bubblefetch "$INSTALL_DIR/"
    sudo ln -sf "$INSTALL_DIR/bubblefetch" "$INSTALL_DIR/bf"
    echo -e "${GREEN}✓ Installed successfully${NC}"
fi

# Create config directory
CONFIG_DIR="$HOME/.config/bubblefetch"
if [ ! -d "$CONFIG_DIR" ]; then
    echo -e "${YELLOW}Creating config directory...${NC}"
    mkdir -p "$CONFIG_DIR"
    echo -e "${GREEN}✓ Config directory created${NC}"
fi

# Copy example config if it doesn't exist
if [ ! -f "$CONFIG_DIR/config.yaml" ]; then
    echo -e "${YELLOW}Copying example config...${NC}"
    cp config.example.yaml "$CONFIG_DIR/config.yaml"
    echo -e "${GREEN}✓ Config file created${NC}"
fi

# Copy themes
THEMES_DIR="$CONFIG_DIR/themes"
if [ ! -d "$THEMES_DIR" ]; then
    echo -e "${YELLOW}Creating themes directory...${NC}"
    mkdir -p "$THEMES_DIR"
fi

echo -e "${YELLOW}Copying themes...${NC}"
cp themes/*.json "$THEMES_DIR/" 2>/dev/null || true
echo -e "${GREEN}✓ Themes installed${NC}"

echo ""
echo -e "${GREEN}╔═══════════════════════════════╗${NC}"
echo -e "${GREEN}║  Installation Complete! 🎉   ║${NC}"
echo -e "${GREEN}╚═══════════════════════════════╝${NC}"
echo ""
echo -e "Run ${BLUE}bubblefetch${NC} or ${BLUE}bf${NC} to get started!"
echo -e "Run ${BLUE}bubblefetch --help${NC} for options"
echo ""
echo -e "Config: ${YELLOW}$CONFIG_DIR/config.yaml${NC}"
echo -e "Themes: ${YELLOW}$THEMES_DIR/${NC}"

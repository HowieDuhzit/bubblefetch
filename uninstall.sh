#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔═══════════════════════════════╗${NC}"
echo -e "${BLUE}║   bubblefetch Uninstaller    ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════╝${NC}"
echo ""

# Remove binary
INSTALL_DIR="/usr/local/bin"
if [ -f "$INSTALL_DIR/bubblefetch" ]; then
    echo -e "${YELLOW}Removing binary...${NC}"
    if [ -w "$INSTALL_DIR" ]; then
        rm "$INSTALL_DIR/bubblefetch"
    else
        sudo rm "$INSTALL_DIR/bubblefetch"
    fi
    echo -e "${GREEN}✓ Binary removed${NC}"
else
    echo -e "${YELLOW}Binary not found in $INSTALL_DIR${NC}"
fi

# Ask about config
read -p "Remove config directory (~/.config/bubblefetch)? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    CONFIG_DIR="$HOME/.config/bubblefetch"
    if [ -d "$CONFIG_DIR" ]; then
        echo -e "${YELLOW}Removing config directory...${NC}"
        rm -rf "$CONFIG_DIR"
        echo -e "${GREEN}✓ Config directory removed${NC}"
    fi
fi

echo ""
echo -e "${GREEN}Uninstallation complete${NC}"

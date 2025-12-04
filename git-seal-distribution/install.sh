#!/bin/bash

# Git-Seal Installation Script
# This script helps you install git-seal on your system

set -e

echo "🔧 Git-Seal Installation Script"
echo "================================"

# Check if running as root (for system-wide installation)
if [ "$EUID" -eq 0 ]; then
    echo "❌ Please run this script as a regular user (not root)"
    echo "   You may be prompted for sudo password later"
    exit 1
fi

# Check if git-seal binary exists
if [ ! -f "./git-seal" ]; then
    echo "❌ git-seal binary not found in current directory"
    echo "   Please ensure you're in the git-seal distribution folder"
    exit 1
fi

# Make it executable
chmod +x git-seal

# Test the binary
echo "✅ Testing git-seal binary..."
./git-seal > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ Binary test passed"
else
    echo "❌ Binary test failed"
    exit 1
fi

# Install to /usr/local/bin
echo "🚀 Installing git-seal to /usr/local/bin..."
if sudo mv git-seal /usr/local/bin/; then
    echo "✅ Successfully installed to /usr/local/bin/git-seal"
else
    echo "❌ Failed to install to /usr/local/bin"
    echo "   You may need to run: sudo mv git-seal /usr/local/bin/"
    exit 1
fi

# Verify installation
echo "🔍 Verifying installation..."
if which git-seal > /dev/null 2>&1; then
    echo "✅ git-seal is now available in your PATH"
    git-seal --version 2>/dev/null || echo "✅ git-seal installed and ready"
else
    echo "⚠️  git-seal installed but may not be in PATH"
    echo "   Try running: export PATH=\$PATH:/usr/local/bin"
fi

echo ""
echo "🎉 Installation Complete!"
echo "========================"
echo "Next steps:"
echo "1. Generate your master key:"
echo "   git-seal keygen"
echo ""
echo "2. Configure Git integration:"
echo "   git-seal setup"
echo ""
echo "3. Read the README.md for usage instructions"
echo ""
echo "⚠️  Remember to backup your ~/.git-seal.key file!"
#!/bin/bash

tools=(
	"github.com/joho/godotenv/cmd/godotenv@v1.5.1"
    "github.com/cosmtrek/air@v1.49.0"
    "github.com/scylladb/gocqlx/v3/cmd/schemagen@latest"
)

# Define colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to display success message
success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Function to display warning message
warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Function to display error message
error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Install tools
install_tool() {
    tool=$1
    go install "$tool"
    if [ $? -eq 0 ]; then
        success "Successfully installed $tool"
    else
        error "Failed to install $tool"
    fi
}

# Main installation script
echo "Installing development tools..."


# Loop through the array and install each tool
for tool in "${tools[@]}"; do
    install_tool "$tool"
done

success "Installation completed."

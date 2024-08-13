#!/bin/bash

tools=(
    "github.com/bufbuild/buf/cmd/buf@v1.29.0"
    "github.com/golang/protobuf/protoc-gen-go@v1.5.3"
    "google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0"
    "github.com/envoyproxy/protoc-gen-validate/cmd/protoc-gen-validate-go@v1.0.4"
    "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.21.0"
    "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.21.0"
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

# rename validate binary since the name is hardcoded in buf
rename_validate_binary() {
   if  mv "./bin/protoc-gen-validate-go" "./bin/protoc-gen-validate"; then
       success "Successfully renamed validate binary"
   else
        error "Failed to rename validate binary"
   fi
}

# Main installation script
echo "Installing development tools..."


# Loop through the array and install each tool
for tool in "${tools[@]}"; do
    install_tool "$tool"
done

success "Installation completed."

rename_validate_binary

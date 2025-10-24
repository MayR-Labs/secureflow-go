#!/bin/bash

# decrypt_secrets.sh
# Decrypts encrypted environment and key files for CI/CD
# Usage: ./decrypt_secrets.sh <ENCRYPTION_PASSWORD>

set -e  # Exit on error

# Load shared config
source "./encrypt_config.sh"
source "./encrypt_functions.sh"

# Check if encryption password is provided
if [ -z "$1" ]; then
  echo -e "${RED}❌ Error: Encryption password not provided${NC}"
  echo "Usage: $0 <ENCRYPTION_PASSWORD>"
  exit 1
fi

ENCRYPTION_PASSWORD="$1"

echo -e "${YELLOW}🔐 Starting decryption process...${NC}"
echo ""

decrypt_file "$OUT_DIR/$ENV_OUT" "$ENV_IN"
decrypt_file "$OUT_DIR/$KEYSTORE_OUT" "$KEYSTORE_IN"
decrypt_file "$OUT_DIR/$KEYPROPERTIES_OUT" "$KEYPROPERTIES_IN"
decrypt_file "$OUT_DIR/$SERVICE_KEY_OUT" "$SERVICE_KEY_IN"

cp "$ENV_IN" ".env"

echo ""
echo ""

echo -e "${GREEN}🎉 All secrets decrypted successfully!${NC}"

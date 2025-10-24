#!/bin/bash

# encrypt_secrets.sh
# Encrypts environment and key files for CI/CD
# Usage: ./encrypt_secrets.sh then enter password when prompted

set -e  # Exit on error

# Load shared config
source "./encrypt_config.sh"
source "./encrypt_functions.sh"


echo -e "${BLUE}🔐 Enter password to encrypt your secrets:${NC}"
read -s ENCRYPTION_PASSWORD

echo ""
echo -e "${BLUE}🔑 (Optional) Enter a password hint (leave blank to skip):${NC}"
read PASSWORD_HINT

echo ""
echo -e "${BLUE}📝 (Optional) Enter a short note (leave blank for default):${NC}"
read SHORT_NOTE

if [ -z "$SHORT_NOTE" ]; then
    SHORT_NOTE="Encrypted secrets for CI/CD"
fi

echo ""

mkdir -p "$OUT_DIR"

{
    echo "Encryption Report"
    echo "================="
    echo ""
    echo "Note: $SHORT_NOTE"
    echo "Password Hint: ${PASSWORD_HINT:-N/A}"
    echo "Created at: $(date)"
    echo "================="
    echo ""
    echo ""
} > "$REPORT_FILE"

encrypt_file "$ENV_IN" "$OUT_DIR/$ENV_OUT"
encrypt_file "$KEYSTORE_IN" "$OUT_DIR/$KEYSTORE_OUT"
encrypt_file "$KEYPROPERTIES_IN" "$OUT_DIR/$KEYPROPERTIES_OUT"
encrypt_file "$SERVICE_KEY_IN" "$OUT_DIR/$SERVICE_KEY_OUT"

echo ""
echo -e "${GREEN}✅ Encryption complete. All files saved to $OUT_DIR${NC}"

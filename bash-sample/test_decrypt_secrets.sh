#!/bin/bash

# test_decrypt_secrets.sh
# Decrypts encrypted environment and key files for CI/CD
# Usage: ./test_decrypt_secrets.sh then enter password when prompted

set -e  # Exit on error

# Load shared config
source "./encrypt_config.sh"
source "./encrypt_functions.sh"


echo -e "${BLUE}🔐 [TEST] - Enter password to encrypt your secrets:${NC}"
read -s ENCRYPTION_PASSWORD

mkdir -p "$TEST_OUT_DIR"

echo -e "${YELLOW}🔐 [TEST] - Starting decryption process...${NC}"
echo ""

decrypt_file "$OUT_DIR/$ENV_OUT" "$TEST_OUT_DIR/.env.prod"
decrypt_file "$OUT_DIR/$KEYSTORE_OUT" "$TEST_OUT_DIR/keystore.jks"
decrypt_file "$OUT_DIR/$KEYPROPERTIES_OUT" "$TEST_OUT_DIR/key.properties"
decrypt_file "$OUT_DIR/$SERVICE_KEY_OUT" "$TEST_OUT_DIR/service-key.json"

echo -e "${GREEN}🎉 All secrets decrypted successfully!${NC}"

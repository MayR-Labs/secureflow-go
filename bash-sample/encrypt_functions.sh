
decrypt_file() {
  local IN_FILE="$1"
  local OUT_FILE="$2"

  echo -e "${YELLOW}📄 Decrypting $IN_FILE...${NC}"

  if openssl enc -aes-256-cbc -d -pbkdf2 \
    -in "$IN_FILE" \
    -out "$OUT_FILE" \
    -k "$ENCRYPTION_PASSWORD"; then

    echo -e "${GREEN}✅ $IN_FILE decrypted successfully -> $OUT_FILE ${NC}"
  else
    echo -e "${RED}❌ Failed to decrypt $IN_FILE${NC}" >&2
    exit 1
  fi

  echo ""
}

encrypt_file() {
    local IN_FILE="$1"
    local OUT_FILE="$2"

    echo -e "${YELLOW}📦 Encrypting $IN_FILE...${NC}"

    if openssl enc -aes-256-cbc -salt -pbkdf2 \
        -in "$IN_FILE" \
        -out "$OUT_FILE" \
        -k "$ENCRYPTION_PASSWORD"; then

        echo -e "${GREEN}✅ $IN_FILE encrypted successfully -> $OUT_FILE${NC}"

        local FILE_SIZE LINE_COUNT LAST_MODIFIED

        FILE_SIZE=$(stat -c%s "$IN_FILE" 2>/dev/null || stat -f%z "$IN_FILE")
        LINE_COUNT=$(wc -l < "$IN_FILE")
        LAST_MODIFIED=$(stat -c%y "$IN_FILE" 2>/dev/null || stat -f"%Sm" -t "%Y-%m-%d %H:%M:%S" "$IN_FILE")

        {
          echo "File:           $IN_FILE"
          echo "Encrypted As:   $OUT_FILE"
          echo "Size (bytes):   $FILE_SIZE"
          echo "Lines:   $LINE_COUNT"
          echo "Last Modified:  $LAST_MODIFIED"
          echo "----------------------------------------"
          echo ""
        } >> "$REPORT_FILE"

    else
        echo -e "${RED}❌ Failed to encrypt $IN_FILE${NC}" >&2
        exit 1
    fi
    echo ""
}

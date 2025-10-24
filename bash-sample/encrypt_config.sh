# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Files and DIRs
OUT_DIR="enc_keys"
TEST_OUT_DIR="test_dec_keys"

ENV_IN=".env.prod"
ENV_OUT=".env.prod.encrypted"

KEYSTORE_IN="android/app/keystore.jks"
KEYSTORE_OUT="keystore.jks.encrypted"

KEYPROPERTIES_IN="android/key.properties"
KEYPROPERTIES_OUT="key.properties.encrypted"

SERVICE_KEY_IN="android/service-key.json"
SERVICE_KEY_OUT="service-key.json.encrypted"

REPORT_FILE="$OUT_DIR/report.txt"

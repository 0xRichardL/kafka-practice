#!/usr/bin/env bash
set -euo pipefail

usage() {
  echo "Usage: $0 <schema-file> <subject> [schema-type] [schema-registry-url]"
  echo "  schema-type: AVRO (default) or PROTOBUF"
  echo "  schema-registry-url: default http://localhost:8081"
  exit 2
}

if [ "${#}" -lt 2 ]; then
  usage
fi

SCHEMA_FILE="$1"
SUBJECT="$2"
SCHEMA_TYPE="${3:-AVRO}"
REGISTRY_URL="${4:-http://localhost:8081}"

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required. Install with 'brew install jq' (Mac)." >&2
  exit 3
fi

if [ ! -f "$SCHEMA_FILE" ]; then
  echo "Schema file not found: $SCHEMA_FILE" >&2
  exit 4
fi

# Read and JSON-escape the schema file
SCHEMA_JSON=$(jq -Rs . < "$SCHEMA_FILE")

PAYLOAD="{\"schema\": $SCHEMA_JSON, \"schemaType\": \"${SCHEMA_TYPE}\"}"

echo "Checking compatibility for subject '$SUBJECT' against $REGISTRY_URL ..."
tmp=$(mktemp)
status=$(curl -sS -o "$tmp" -w "%{http_code}" \
  -X POST "${REGISTRY_URL%/}/compatibility/subjects/${SUBJECT}/versions" \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  -d "$PAYLOAD" || true)

body=$(cat "$tmp"); rm -f "$tmp"

if [ "$status" = "200" ]; then
  is_compatible=$(echo "$body" | jq -r '.is_compatible // empty')
  if [ "$is_compatible" = "true" ]; then
    echo "Compatibility check passed."
  else
    echo "Compatibility check failed. Response:"
    echo "$body" | jq .
    exit 5
  fi
else
  echo "Compatibility check returned HTTP $status. Proceeding to register (subject may be new)."
  [ -n "$body" ] && echo "$body" | jq . || true
fi

echo "Registering schema to subject '$SUBJECT' ..."
resp=$(curl -sS -X POST "${REGISTRY_URL%/}/subjects/${SUBJECT}/versions" \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  -d "$PAYLOAD")

echo "Register response:"
echo "$resp" | jq .
#!/usr/bin/env bash
set -euo pipefail
SEARCH_INDEX_HTTP_ADDR="${SEARCH_INDEX_HTTP_ADDR:-:8083}" go run ./cmd/search-index-service

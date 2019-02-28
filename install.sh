#!/bin/sh
set -e

OUTPUT_FILE='/usr/local/bin/krontab'
RELEASES_URL="https://github.com/jacobtomlinson/krontab/releases"

last_version() {
  curl -sL -o /dev/null -w %{url_effective} "$RELEASES_URL/latest" |
    rev |
    cut -f1 -d'/'|
    rev
}

download() {
  test -z "$VERSION" && VERSION="$(last_version)"
  test -z "$VERSION" && {
    echo "Unable to get goreleaser version." >&2
    exit 1
  }
  rm -f "$OUTPUT_FILE"
  curl -s -L -o "$OUTPUT_FILE" \
    "$RELEASES_URL/download/$VERSION/krontab-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)"
}

echo "Downloading krontab..."
download
chmod +x $OUTPUT_FILE
echo "Done!"

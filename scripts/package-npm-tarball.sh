#!/usr/bin/env bash

# Build cross-platform binaries and create an npm-compatible tarball ready to upload to GitHub Releases.
# Usage: scripts/package-npm-tarball.sh <version>
# Example: scripts/package-npm-tarball.sh 0.1.0

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VERSION="${1:-${VERSION:-}}"

if [[ -z "$VERSION" ]]; then
  echo "Usage: $0 <version>" >&2
  exit 1
fi

STAGE_DIR="$ROOT/.artifacts/npm-package"
DIST_DIR="$ROOT/dist/npm"

rm -rf "$STAGE_DIR"
mkdir -p "$STAGE_DIR/bin/linux-amd64" "$STAGE_DIR/bin/windows-amd64" "$DIST_DIR"

echo "Building linux/amd64 binary..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "$STAGE_DIR/bin/linux-amd64/magnet2torrent" ./cmd/magnet2torrent
chmod +x "$STAGE_DIR/bin/linux-amd64/magnet2torrent"

echo "Building windows/amd64 binary..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o "$STAGE_DIR/bin/windows-amd64/magnet2torrent.exe" ./cmd/magnet2torrent

cat > "$STAGE_DIR/bin/run.js" <<'EOF'
#!/usr/bin/env node
const { platform, arch, argv, env, exit } = process;
const path = require('path');
const { spawnSync } = require('child_process');

const targets = {
  linux: { x64: 'bin/linux-amd64/magnet2torrent' },
  win32: { x64: 'bin/windows-amd64/magnet2torrent.exe' },
};

const targetRel = targets[platform]?.[arch];
if (!targetRel) {
  console.error(`magnet2torrent: no binary available for ${platform}/${arch}`);
  exit(1);
}

const binaryPath = path.join(__dirname, targetRel);
const result = spawnSync(binaryPath, argv.slice(2), { stdio: 'inherit', env });
if (result.error) {
  console.error(result.error.message);
  exit(1);
}

exit(result.status ?? 1);
EOF

chmod +x "$STAGE_DIR/bin/run.js"

cat > "$STAGE_DIR/package.json" <<EOF
{
  "name": "magnet2torrent",
  "version": "$VERSION",
  "description": "Magnet link to torrent CLI packaged for npm installs",
  "bin": {
    "magnet2torrent": "./bin/run.js"
  },
  "files": [
    "bin/"
  ],
  "license": "UNLICENSED",
  "os": [
    "linux",
    "win32"
  ],
  "cpu": [
    "x64"
  ]
}
EOF

cat > "$STAGE_DIR/README.md" <<EOF
# magnet2torrent npm package

This tarball packages the magnet2torrent CLI for npm installation.

## Installation

\`\`\`bash
npm install -g https://github.com/you/repo/releases/download/v$VERSION/magnet2torrent-$VERSION.tgz
\`\`\`
EOF

pushd "$ROOT" >/dev/null
tarball_name=$(npm pack "$STAGE_DIR" | tail -n1)
popd >/dev/null

mv "$ROOT/$tarball_name" "$DIST_DIR/"

echo "Created npm package: $DIST_DIR/$tarball_name"
echo "Upload this tarball to a GitHub Release at: https://github.com/you/repo/releases/tag/v$VERSION"

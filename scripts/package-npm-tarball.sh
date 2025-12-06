#!/usr/bin/env bash

# Build cross-platform binaries and create an npm-compatible tarball ready to upload to GitHub Releases.
# Usage: scripts/package-npm-tarball.sh <version>
# Example: scripts/package-npm-tarball.sh 0.1.0

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VERSION="${1:-${VERSION:-}}"
REPO_OWNER="${REPO_OWNER:-alacasse}"
REPO_NAME="${REPO_NAME:-magnet2torrent}"

if [[ -z "$VERSION" ]]; then
  echo "Usage: $0 <version>" >&2
  exit 1
fi

STAGE_DIR="$ROOT/.artifacts/npm-package"
DIST_DIR="$ROOT/dist/npm"

rm -rf "$STAGE_DIR"
mkdir -p "$STAGE_DIR/bin/linux-amd64" "$STAGE_DIR/bin/windows-amd64" "$STAGE_DIR/scripts" "$DIST_DIR"

echo "Building linux/amd64 binary..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "$STAGE_DIR/bin/linux-amd64/magnet2torrent" ./cmd/magnet2torrent
chmod +x "$STAGE_DIR/bin/linux-amd64/magnet2torrent"

echo "Building windows/amd64 binary..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o "$STAGE_DIR/bin/windows-amd64/magnet2torrent.exe" ./cmd/magnet2torrent

echo "Copying helper scripts..."
cp "$ROOT/scripts/register-magnet-linux.sh" "$STAGE_DIR/scripts/"
cp "$ROOT/scripts/register-magnet-windows.ps1" "$STAGE_DIR/scripts/"
chmod +x "$STAGE_DIR/scripts/register-magnet-linux.sh"

cat > "$STAGE_DIR/bin/run.js" <<'EOF'
#!/usr/bin/env node
const { platform, arch, argv, env, exit } = process;
const path = require('path');
const fs = require('fs');
const { spawnSync } = require('child_process');

const targets = {
  linux: { x64: 'linux-amd64/magnet2torrent' },
  win32: { x64: 'windows-amd64/magnet2torrent.exe' },
};

const targetRel = targets[platform]?.[arch];
if (!targetRel) {
  console.error(`magnet2torrent: no binary available for ${platform}/${arch}`);
  exit(1);
}

const binaryPath = path.join(__dirname, targetRel);

if (!fs.existsSync(binaryPath)) {
  console.error(`magnet2torrent: binary missing at ${binaryPath}. Reinstall the package or re-run the installer.`);
  exit(1);
}

const result = spawnSync(binaryPath, argv.slice(2), { stdio: 'inherit', env });
if (result.error) {
  console.error(`magnet2torrent: failed to launch binary at ${binaryPath}: ${result.error.message}`);
  exit(1);
}

exit(result.status ?? 0);
EOF

chmod +x "$STAGE_DIR/bin/run.js"

cat > "$STAGE_DIR/postinstall.js" <<'EOF'
#!/usr/bin/env node
const { platform, env, exit } = process;
const path = require('path');
const fs = require('fs');
const { spawnSync } = require('child_process');

const skip = env.REGISTER_MAGNET === '0';
if (skip) {
  console.log('magnet2torrent: skipping magnet handler registration (REGISTER_MAGNET=0)');
  exit(0);
}

const dryRun = env.DRY_RUN === '1';
const scriptsDir = path.join(__dirname, 'scripts');

function run(cmd, args) {
  const res = spawnSync(cmd, args, { stdio: 'inherit' });
  if (res.error) {
    console.error(`magnet2torrent: failed to run ${cmd}: ${res.error.message}`);
    return false;
  }
  if (res.status !== 0) {
    console.error(`magnet2torrent: registration command exited with ${res.status}`);
    return false;
  }
  return true;
}

function ensureScript(name) {
  const full = path.join(scriptsDir, name);
  if (!fs.existsSync(full)) {
    console.error(`magnet2torrent: registration script missing: ${full}`);
    return null;
  }
  return full;
}

let ok = true;
if (platform === 'linux') {
  const script = ensureScript('register-magnet-linux.sh');
  if (script) {
    ok = run(script, dryRun ? ['--dry-run'] : []);
  } else {
    ok = false;
  }
} else if (platform === 'win32') {
  const script = ensureScript('register-magnet-windows.ps1');
  if (script) {
    const args = ['-ExecutionPolicy', 'Bypass', '-File', script];
    if (dryRun) { args.push('-DryRun'); }
    ok = run('pwsh', args);
  } else {
    ok = false;
  }
} else {
  console.log(`magnet2torrent: postinstall registration not supported on ${platform}`);
}

if (!ok) {
  console.error('magnet2torrent: magnet handler registration failed; run the platform script manually later.');
}
EOF

cat > "$STAGE_DIR/package.json" <<EOF
{
  "name": "magnet2torrent",
  "version": "$VERSION",
  "description": "Magnet link to torrent CLI packaged for npm installs",
  "bin": {
    "magnet2torrent": "./bin/run.js"
  },
  "files": [
    "bin/",
    "scripts/",
    "postinstall.js"
  ],
  "scripts": {
    "postinstall": "node postinstall.js"
  },
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
npm install -g https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/v$VERSION/magnet2torrent-$VERSION.tgz
\`\`\`
EOF

pushd "$ROOT" >/dev/null
tarball_name=$(npm pack "$STAGE_DIR" | tail -n1)
popd >/dev/null

mv "$ROOT/$tarball_name" "$DIST_DIR/"

echo "Created npm package: $DIST_DIR/$tarball_name"
echo "Upload this tarball to a GitHub Release at: https://github.com/$REPO_OWNER/$REPO_NAME/releases/tag/v$VERSION"

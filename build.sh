#!/bin/bash

set -e

echo "🎮 Building LBaaS Packet Catcher (Modular Power-up Edition)..."

# Clean previous builds
rm -rf dist
mkdir -p dist

# List of OS/ARCH combinations
platforms=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

for platform in "${platforms[@]}"; do
  os="${platform%%/*}"
  arch="${platform##*/}"
  output_name="lbbaspack"
  ext=""
  
  if [ "$os" = "windows" ]; then
    output_name+=".exe"
    ext=".zip"
  else
    ext=".tar.gz"
  fi

  output_dir="dist/lbbaspack_${os}_${arch}"
  mkdir -p "$output_dir"

  echo "🔧 Building for $os/$arch..."

  GOOS="$os" GOARCH="$arch" go build -o "$output_dir/$output_name" .

  if [ "$ext" = ".zip" ]; then
    (cd "$output_dir" && zip "../lbbaspack_${os}_${arch}$ext" "$output_name")
  else
    (cd "$output_dir" && tar -czf "../lbbaspack_${os}_${arch}$ext" "$output_name")
  fi

  rm -rf "$output_dir"
done

echo "✅ All builds completed and packaged into dist/"


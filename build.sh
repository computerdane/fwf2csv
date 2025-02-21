#!/usr/bin/env bash

cd lib
if ! go test; then
  echo "Tests failed! Exiting..."
  exit 1
fi
cd ..

if [ -z "${VERSION}" ]; then
  VERSION=devel
fi

rm -rf dist
mkdir -p dist

os_archs=("darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64")
for os_arch in "${os_archs[@]}"; do
  IFS="/" read -r os arch <<< "$os_arch"
  output_name="fwf2csv-${os}-${arch}"
  echo "Building for $os/$arch..."
  GOOS=$os GOARCH=$arch go build -o "dist/$output_name" \
    -ldflags="-X main.Version=$VERSION"
done


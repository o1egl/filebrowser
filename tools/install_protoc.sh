#!/bin/bash

set -euf -o pipefail

readonly version=3.15.6

case "$(uname -s)" in
  Linux*)
    filename=protoc-${version}-linux-x86_64.zip
  ;;
  Darwin*)
    filename=protoc-${version}-osx-x86_64.zip
  ;;
  CYGWIN*|MINGW*)
    filename=protoc-${version}-win64.zip
  ;;
  *)
    echo "Unsupported system"
    exit 1
  ;;
esac

rm -rf protoc
curl -L "https://github.com/protocolbuffers/protobuf/releases/download/v${version}/${filename}" -o protoc.zip
unzip -o protoc.zip -d protoc
mkdir -p bin/
mv protoc/bin/protoc bin/protoc
rm protoc.zip

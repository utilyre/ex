#!/bin/sh

rm -rf build
mkdir -p build

cp -r assets build/assets
cp -r views build/views
go build -o build/server

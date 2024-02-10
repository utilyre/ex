#!/bin/sh

rm -rf build
mkdir -p build

cp -r public build/public
cp -r views build/views
go build -o build/server

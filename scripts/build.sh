#!/bin/bash

if [ ! -d ./bin ]; then
    mkdir ./bin
fi

echo "Build Low-level Container Runtime..."
sh ./scripts/futaba/build.sh
echo "Done."

echo "Build High-level Container Runtime..."
sh ./scripts/hitoha/build.sh
echo "Done."

echo "Build Karakuri CLI..."
sh ./scripts/karakuri/build.sh
echo "Done."


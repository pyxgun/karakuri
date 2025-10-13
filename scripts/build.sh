#!/bin/bash

if [ ! -d ./bin ]; then
    mkdir ./bin
fi

echo -n "Build Low-level Container Runtime..."
sh ./scripts/futaba/build.sh
echo "Done."

echo -n "Build High-level Container Runtime..."
sh ./scripts/hitoha/build.sh
echo "Done."

echo -n "Build Karakuri CLI..."
sh ./scripts/karakuri/build.sh
echo "Done."

sudo systemctl stop hitoha.service
sudo cp ./bin/* /bin

sudo cp ./scripts/hitoha/hitoha.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable hitoha.service

sudo systemctl restart hitoha.service
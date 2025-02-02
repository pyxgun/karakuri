#!/bin/bash

cmd_dir="./cmd/hitoha"
cd $cmd_dir

outdir="../../bin"
outname="hitoha"

sudo systemctl stop hitoha.service

go build -o $outdir/$outname
sudo cp $outdir/$outname /bin

sudo cp ../../scripts/hitoha/hitoha.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable hitoha.service

sudo systemctl restart hitoha.service
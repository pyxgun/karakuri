#!/bin/bash

cmd_dir="./cmd/hitoha"
cd $cmd_dir

outdir="../../bin"
outname="hitoha"

go build -o $outdir/$outname

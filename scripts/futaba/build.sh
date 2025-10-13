#!/bin/bash

cmd_dir="./cmd/futaba"
cd $cmd_dir

outdir="../../bin"
outname="futaba"

go build -o $outdir/$outname
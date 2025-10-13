#!/bin/bash

cmd_dir="./cmd/karakuri"
cd $cmd_dir

outdir="../../bin"
outname="karakuri"

go build -o $outdir/$outname

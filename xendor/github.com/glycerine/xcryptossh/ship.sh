#!/bin/bash
where=~/go/src/github.com/glycerine/xcryptossh
cp *.go ${where}
cp README.md ${where}
cd ${where}; ./fix.sh

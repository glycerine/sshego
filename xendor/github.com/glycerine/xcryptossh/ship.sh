#!/bin/bash
where=${GOPATH}/src/github.com/glycerine/xcryptossh
cp *.go ${where}
for i in agent knownhosts terminal test testdata; do
  cp $i/* ${where}/$i/
done
cp README.md ${where}
cd ${where}; ./fix.sh

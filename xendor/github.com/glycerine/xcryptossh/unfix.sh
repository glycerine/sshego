#!/bin/bash

for i in `ls -1 *test.go`; do
  perl -pi -e 's$\"github.com/glycerine/xcryptossh/$\"github.com/glycerine/sshego/xendor/github.com/glycerine/xcryptossh/$g' $i
done

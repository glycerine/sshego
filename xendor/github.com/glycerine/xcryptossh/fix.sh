#!/bin/bash
perl -pi -e 's$github.com/glycerine/sshego/xendor/$$g' *.go
for i in agent knownhosts terminal test testdata; do
  perl -pi -e 's$github.com/glycerine/sshego/xendor/$$g' $i/*.go
done

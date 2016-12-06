.phony: all

all:
	go install github.com/glycerine/sshego
	go install github.com/glycerine/sshego/cmd/gosshtun


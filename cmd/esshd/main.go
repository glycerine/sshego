package main

import (
	"flag"
	"log"
	"os"

	tun "github.com/glycerine/gosshtun"
)

const ProgramName = "gosshtun"

func main() {

	myflags := flag.NewFlagSet(ProgramName, flag.ExitOnError)
	cfg := tun.NewGosshtunConfig()
	cfg.DefineFlags(myflags)
	err := myflags.Parse(os.Args[1:])
	err = cfg.ValidateConfig()
	if err != nil {
		log.Fatalf("%s command line flag error: '%s'", ProgramName, err)
	}

	cfg.NewEsshd()
	cfg.Esshd.Start()
	select {}
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

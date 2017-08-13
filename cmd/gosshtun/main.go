package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	tun "github.com/glycerine/sshego"
)

const ProgramName = "gosshtun"

func main() {

	myflags := flag.NewFlagSet(ProgramName, flag.ExitOnError)
	cfg := tun.NewSshegoConfig()
	cfg.DefineFlags(myflags)
	err := myflags.Parse(os.Args[1:])
	if cfg.ShowVersion {
		fmt.Printf("\n%v\n", tun.SourceVersion())
		os.Exit(0)
	}
	err = cfg.ValidateConfig()
	if err != nil {
		log.Fatalf("%s command line flag error: '%s'", ProgramName, err)
	}
	//p("cfg = %#v", cfg)
	h, err := tun.NewKnownHosts(cfg.ClientKnownHostsPath, tun.KHJson)
	panicOn(err)
	cfg.KnownHosts = h

	if cfg.WriteConfigOut != "" {
		var o io.WriteCloser
		if cfg.WriteConfigOut == "-" {
			o = os.Stdout
		} else {
			o, err = os.Create(cfg.WriteConfigOut)
			if err != nil {
				panic(err)
			}
		}
		err = cfg.SaveConfig(o)
		if err != nil {
			panic(err)
		}
	}

	if cfg.AddUser != "" {
		tun.AddUserAndExit(cfg)
	}

	if cfg.DelUser != "" {
		tun.DelUserAndExit(cfg)
	}

	passphrase := ""
	totpUrl := ""
	ctx := context.Background()

	_, _, err = cfg.SSHConnect(h, cfg.Username, cfg.PrivateKeyPath,
		cfg.SSHdServer.Host, cfg.SSHdServer.Port, passphrase, totpUrl, ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	if err != nil {
		panic(err)
	}
	if !cfg.WriteConfigOnly {
		select {}
	}
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}

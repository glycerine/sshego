package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	tun "github.com/glycerine/sshego"
	"github.com/skratchdot/open-golang/open"
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
		addUserAndExit(cfg)
	}

	if cfg.DelUser != "" {
		delUserAndExit(cfg)
	}

	passphrase := ""
	totpUrl := ""

	_, err = cfg.SSHConnect(h, cfg.Username, cfg.PrivateKeyPath,
		cfg.SSHdServer.Host, cfg.SSHdServer.Port, passphrase, totpUrl)
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

func addUserAndExit(cfg *tun.SshegoConfig) {

	err := cfg.NewHostDb()
	panicOn(err)

	mylogin := cfg.AddUser

	if cfg.HostDb.UserExists(mylogin) {
		fmt.Fprintf(os.Stderr, "\nerror: user '%s' already exists. If you want to replace them, use -deluser first.\n", mylogin)
		os.Exit(1)
	}
	var ok bool
	ok, err = cfg.HostDb.ValidLogin(mylogin)
	if !ok {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	var myemail string
	fmt.Printf("\nEnter the email address for '%s' (for backups/recovery): ",
		mylogin)
	myemail, err = reader.ReadString('\n')
	panicOn(err)
	myemail = strings.Trim(myemail, "\n\r\t ")
	ok, err = cfg.HostDb.ValidEmail(myemail)
	if !ok {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n'%s' should be a valid email, but may be fake if you insist.\n\nIf this email is real, we can backup the passphrase you will create in the next step. This convenience feature is optional, but strongly recommended.  The RSA key and One-time-password secret will not be sent by email, but their locations on disk will be noted for ease of retreival.\n\nDo you want to backup the passphase to email '%s'? [y/n]:", myemail, myemail)
	yn, err := reader.ReadString('\n')
	panicOn(err)
	yn = strings.ToLower(strings.Trim(yn, "\n\r\t "))
	sendEmail := false
	if yn == "y" || yn == "yes" {
		sendEmail = true
		fmt.Printf("\n\n Very good. Upon completion, your passphrase "+
			"will backed up to '%s'\n\n", myemail)
	} else {
		fmt.Printf("\n\n As you wish. No email will be sent.\n\n")
	}

	var fullname string
	fmt.Printf("\nCorresponding to '%s'/'%s', enter the first and last name (e.g. 'John Q. Smith'). This helps identify the account during maintenance. First and last name:\n", mylogin, myemail)
	fullname, err = reader.ReadString('\n')
	panicOn(err)
	fullname = strings.Trim(fullname, "\n\r\t ")

	var pw string
	if !cfg.SkipPassphrase {
		pw, err = tun.PromptForPassword(cfg.AddUser)
		if err != nil {
			fmt.Printf("\n%v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\n account: '%s', passphrase: '%s'\n", cfg.AddUser, pw)
	}
	if !cfg.SkipRSA {
		fmt.Printf("\n generating a strong RSA key, this may take a 5-10 seconds... \n")
	}
	user := tun.NewUser()
	user.MyLogin = mylogin
	user.MyEmail = myemail
	user.MyFullname = fullname
	user.ClearPw = pw
	user.Issuer = "gosshtun"

	var toptPath, qrPath, rsaPath string

	var prt tun.TcpPort
	prt.Port = cfg.SshegoSystemMutexPort

	limitMsec := 5000
	err = prt.Lock(limitMsec)
	if err == tun.ErrCouldNotAquirePort {
		// already running...
		p("we see gosshtun is already running and has the xport open")
		toptPath, qrPath, rsaPath, err = cfg.TcpClientUserAdd(user)
	} else {
		p("we got xport, so while holding it, modify the database directly")
		// we must do it ourselves; other process is not
		// up and we now hold the port (listening on it) as a lock.
		toptPath, qrPath, rsaPath, err = cfg.HostDb.AddUser(
			mylogin, myemail, pw, "gosshtun", fullname)
		prt.Unlock()
	}
	if err != nil {
		es := err.Error()
		if strings.HasPrefix(es,
			"bad email: give a full email address.") {
			fmt.Printf("\n%s\n", es)
			os.Exit(1)
		}
	}
	panicOn(err)

	hostname, _ := os.Hostname()
	var plain bytes.Buffer
	var html bytes.Buffer
	both := io.MultiWriter(&plain, &html)
	ip := tun.GetExternalIP()
	userEnv := os.Getenv("USER")

	fmt.Fprintf(&html, "<html><body>")
	fmt.Fprintf(both, "%s:\n\n", fullname)
	fmt.Fprintf(&html, "<p><font size=3>")

	fmt.Fprintf(both, "## ===============================================\n")
	fmt.Fprintf(&html, "<br>")
	fmt.Fprintf(both, "## \n")
	fmt.Fprintf(&html, "<br>")
	authDetail := "tri-factor"
	if cfg.SkipPassphrase || cfg.SkipRSA || cfg.SkipTOTP {
		authDetail = ""
	}
	authList := cfg.GenAuthString()
	fmt.Fprintf(both, "##  %s auth details:\n", authDetail)
	fmt.Fprintf(&html, "<br>")
	fmt.Fprintf(both, "## \n")
	fmt.Fprintf(&html, "<br>")
	fmt.Fprintf(both, "##    %s\n", authList)
	fmt.Fprintf(&html, "<br>")
	fmt.Fprintf(both, "## \n")
	fmt.Fprintf(&html, "<br>")
	fmt.Fprintf(both, "## ==  run: gosshtun -adduser %s\n",
		cfg.AddUser)
	fmt.Fprintf(&html, "<br>")
	fmt.Fprintf(both, "## ==  run on host:  %s (%s)\n", hostname, ip)
	fmt.Fprintf(&html, "<br>")
	fmt.Fprintf(both, "## ==  run by user:  %s\n", userEnv)
	fmt.Fprintf(&html, "<br>")
	fmt.Fprintf(both, "## ==  at UTC time:  %s\n", time.Now().UTC())
	fmt.Fprintf(&html, "<br>")
	fmt.Fprintf(both, "## ===============================================\n")
	fmt.Fprintf(&html, "<br>")

	fmt.Fprintf(&plain, "\nyour login:\n%s\n\n", mylogin)
	fmt.Fprintf(&html, "\n<p>your login:<br>\n<b>%s</b><p>\n\n", mylogin)

	fmt.Fprintf(&plain, "your email:\n%s\n\n", myemail)
	fmt.Fprintf(&html, "\n<p>your email:<br>\n<b>%s</b><p>\n\n", myemail)

	fmt.Fprintf(&plain, "your fullname:\n%s\n\n", fullname)
	fmt.Fprintf(&html, "\n<p>your fullname:<br>\n<b>%s</b><p>\n\n", fullname)

	if !cfg.SkipPassphrase {
		fmt.Fprintf(&plain, "Passphrase:\n%v\n\n", pw)
		fmt.Fprintf(&html, "<p>Passphrase:\n<br>\n<b>%v</b>\n\n", pw)
		fmt.Fprintf(&html, "<p><p>")
	}
	if !cfg.SkipTOTP {
		fmt.Fprintf(&plain, "GoogleAuthenticator (time-based-one-time-password; RFC6238) secret location (on host %s):\n%s\n\n", hostname, toptPath)
		fmt.Fprintf(&html, "GoogleAuthenticator (time-based-one-time-password; RFC6238) secret location (on host %s):\n<br>\n<b>%s</b><p><p>\n", hostname, toptPath)

		qrUrl := fmt.Sprintf("file://%s", qrPath)

		fmt.Printf("\n checking if we should open the QR-code automajically...\n")
		if runtime.GOOS == "darwin" { // "windows", "linux"
			fmt.Printf("...runtime.GOOS='%s'; try to open the QR-code\n",
				runtime.GOOS)
			open.Start(qrUrl)
		}

		fmt.Fprintf(&plain, "GoogleAuthenticator QR-code url (on host %s):\n%s\n\n", hostname, qrPath)
		fmt.Fprintf(&html, "GoogleAuthenticator QR-code url (on host %s):\n<br><b><a href=\"%s\" target=\"_blank\">%s</a></b><p>\n\n", hostname, qrUrl, qrUrl)
		fmt.Fprintf(&html, "<p>")
	}
	if !cfg.SkipRSA {
		fmt.Fprintf(&plain, "Your new RSA Private key is here (on host %s):\n%s\n\n", hostname, rsaPath)
		fmt.Fprintf(&html, "Your new RSA Private key is here (on host %s):\n<br><b>%s\n\n</b><p>", hostname, rsaPath)

		fmt.Fprintf(&plain, "Your new RSA Public key is here (on host %s):\n%s\n\n", hostname, rsaPath+".pub")
		fmt.Fprintf(&html, "Your new RSA Public key is here (on host %s):\n<br><b>%s</b><p>\n\n", hostname, rsaPath+".pub")
	}
	fmt.Fprintf(&html, "</body></html>")

	os.Stdout.Write(plain.Bytes())

	if sendEmail {
		if cfg.MailCfg.Domain == "" ||
			cfg.MailCfg.PublicApiKey == "" ||
			cfg.MailCfg.SecretApiKey == "" {
			fmt.Printf("\n\n alert! -- mailgun not configured; not sending backup email.\n\n")
		} else {
			subject := fmt.Sprintf("gosshtun passphrase backup "+
				"- from %s@%s (%s)", userEnv, hostname, ip)
			senderEmail := fmt.Sprintf("%s@%s", userEnv, hostname)

			pl := string(plain.Bytes())
			ht := string(html.Bytes())

			id, err := cfg.MailCfg.SendEmail(senderEmail,
				subject,
				pl,
				ht,
				myemail)
			if err != nil {
				fmt.Printf("\n mailCfg.SendEmail() failed: %s\n", err)
			} else {
				fmt.Printf("\n backup email sent; id: '%s'\n", id)
			}
		}
	}

	os.Exit(0)
}

func delUserAndExit(cfg *tun.SshegoConfig) {

	fmt.Printf("\ndeleting user '%s'...\n", cfg.DelUser)

	user := tun.NewUser()
	user.MyLogin = cfg.DelUser

	var prt tun.TcpPort
	prt.Port = cfg.SshegoSystemMutexPort

	limitMsec := 5000
	err := prt.Lock(limitMsec)
	if err == tun.ErrCouldNotAquirePort {
		// already running...
		p("we see gosshtun is already running and has the xport open")
		err = cfg.TcpClientUserDel(user)
		if err != nil {
			fmt.Printf("\n error: %s\n", err)
			os.Exit(1)
		}
	} else {
		p("we got xport, so while holding it, modify the database directly")
		// we must do it ourselves; other process is not
		// up and we now hold the port (listening on it) as a lock.

		err := cfg.NewHostDb()
		if err != nil {
			fmt.Printf("\n%s\n", err)
			os.Exit(1)
		}
		err = cfg.HostDb.DelUser(cfg.DelUser)
		if err != nil {
			fmt.Printf("\n%s\n", err)
			os.Exit(1)
		}
		prt.Unlock()
	}
	fmt.Printf("\n deleted user '%s'\n", cfg.DelUser)
	os.Exit(0)
}

package gosshtun

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mailgun/mailgun-go"
	//"github.com/shurcooL/go-goon"
	"io"
	"os"
	"strings"
)

// MailgunConfig sets up sending
// backup emails through Mailgun.
// See https://mailgun.com.
//
type MailgunConfig struct {

	// MAILGUN_DOMAIN
	Domain string

	// MAILGUN_PUBLIC_API_KEY
	PublicApiKey string

	//MAILGUN_SECRET_API_KEY
	SecretApiKey string
}

// LoadConfig reads configuration from a file, expecting
// KEY=value pair on each line;
// values optionally enclosed in double quotes.
func (c *MailgunConfig) LoadConfig(path string) error {
	if !fileExists(path) {
		return fmt.Errorf("path '%s' does not exist", path)
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	bufIn := bufio.NewReader(file)
	lineNum := int64(1)
	for {
		lastLine, err := bufIn.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF && len(lastLine) == 0 {
			break
		}
		line := string(lastLine)
		line = strings.Trim(line, "\n\r\t ")

		if len(line) > 0 && line[0] == '#' {
			// comment, ignore
		} else {

			splt := strings.SplitN(line, "=", 2)
			if len(splt) != 2 {
				/*fmt.Fprintf(os.Stderr, "ignoring malformed (path: '%s') "+
				"config line(%v): '%s'\n",
				path, lineNum, line)
				*/
				continue
			}
			key := strings.Trim(splt[0], "\t\n\r ")
			val := strings.Trim(splt[1], "\t\n\r ")

			val = trim(val)

			switch key {
			case "MAILGUN_DOMAIN":
				c.Domain = val

			case "MAILGUN_PUBLIC_API_KEY":
				c.PublicApiKey = val

			case "MAILGUN_SECRET_API_KEY":
				c.SecretApiKey = val
			}
		}
		lineNum++

		if err == io.EOF {
			break
		}
	}

	return nil
}

// SaveConfig writes the config structs to the given io.Writer
func (c *MailgunConfig) SaveConfig(fd io.Writer) error {

	_, err := fmt.Fprintf(fd, `#
# config for Mailgun:
#
`)
	if err != nil {
		return err
	}
	fmt.Fprintf(fd, "MAILGUN_DOMAIN=\"%s\"\n", c.Domain)
	fmt.Fprintf(fd, "MAILGUN_PUBLIC_API_KEY=\"%s\"\n", c.PublicApiKey)
	fmt.Fprintf(fd, "MAILGUN_SECRET_API_KEY=\"%s\"\n", c.SecretApiKey)
	return nil
}

// DefineFlags should be called before myflags.Parse().
func (c *MailgunConfig) DefineFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.Domain, "mailgun-domain", "", "(supports -adduser) mailgun domain from which to send. Only required if sending backup emails.")
	fs.StringVar(&c.PublicApiKey, "mailgun-pubkey", "", "(supports -adduser) mailgun public api key. Only required if sending backup emails.")
	fs.StringVar(&c.SecretApiKey, "mailgun-secretkey", "", "(supports -adduser) mailgun secret api key. Only required if sending backup emails.")
}

// ValidateConfig should be called after myflags.Parse().
func (c *MailgunConfig) ValidateConfig() error {
	return nil
}

func (c *MailgunConfig) SendEmail(senderEmail, subject, plain, html, recipEmail string) (string, error) {
	/*
		mg := mailgun.NewMailgun(c.Domain, c.SecretApiKey, c.PublicApiKey)
		m := mailgun.NewMessage(senderEmail, subject, body, recipEmail)
		m.SetHtml(body)
	*/

	mg := mailgun.NewMailgun(c.Domain, c.SecretApiKey, c.PublicApiKey)
	from := senderEmail

	m := mg.NewMessage(from, subject, plain)
	m.SetHtml(html)
	m.AddRecipient(recipEmail)

	_, id, err := mg.Send(m)
	//fmt.Printf("debug: details of messagex:\n")
	//goon.Dump(m)
	return id, err
}

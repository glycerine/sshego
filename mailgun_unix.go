// +build darwin linux
// +build !windows,!nacl,!plan9

package sshego

import (
	"context"
	"github.com/mailgun/mailgun-go"
	//"github.com/shurcooL/go-goon"
)

func (c *MailgunConfig) SendEmail(senderEmail, subject, plain, html, recipEmail string) (string, error) {
	/*
		mg := mailgun.NewMailgun(c.Domain, c.SecretApiKey, c.PublicApiKey)
		m := mailgun.NewMessage(senderEmail, subject, body, recipEmail)
		m.SetHtml(body)
	*/

	mg := mailgun.NewMailgun(c.Domain, c.SecretApiKey)
	from := senderEmail

	m := mg.NewMessage(from, subject, plain)
	m.SetHtml(html)
	m.AddRecipient(recipEmail)

	_, id, err := mg.Send(context.Background(), m)
	return id, err
}

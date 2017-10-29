// +build !darwin,!linux,
// +build windows,nacl,plan9

package sshego

func (c *MailgunConfig) SendEmail(senderEmail, subject, plain, html, recipEmail string) (string, error) {
	panic("not implimented")
}

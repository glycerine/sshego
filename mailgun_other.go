// +build windows
// +build !darwin,!linux,

package sshego

func (c *MailgunConfig) SendEmail(senderEmail, subject, plain, html, recipEmail string) (string, error) {
	panic("not implimented")
	return "", nil
}

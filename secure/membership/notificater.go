// Copyright 2012 The Toys Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package membership

import (
	"bytes"
	"net/smtp"
	"text/template"
)

type Notificater interface {
	AccountAdded(User) error
	PasswordChanged(User) error
	AccountInfoChanged(User) error
	AccountPrivilegeChanged(User) error
}

type SMTPNotificater struct {
	Auth smtp.Auth
	// Outgoing Mail (SMTP) Server Address, eg: smtp.gmail.com
	Addr string
	// Port of SMTP Server, eg: 587
	Port int
	// Email use to send, eg: no-reply@gmail.com
	Email string
	//
	accAddedTmpl    *template.Template
	passChangedTmpl *template.Template
	infoChangedTmpl *template.Template
	privChangedTmpl *template.Template
}

func NewSMTPNotificater(email, pass, addr string, port int) *SMTPNotificater {
	n := &SMTPNotificater{}
	n.Addr = addr
	n.Port = port
	n.Email = email
	n.Auth = smtp.PlainAuth("", n.Email, pass, n.Addr)
	n.accAddedTmpl = template.Must(template.New("accAddedTmpl").Parse(`
		Subject: noreply: Your new Account.

		Hi {{.GetEmail}},
		You just reg an account with us!
		{{if not .IsApproved}}
		But you need to confirm by this code:
		{{.GetConfirmCode}}
		{{end}}
		Thanks!
	`))
	n.passChangedTmpl = template.Must(template.New("passChangedTmpl").Parse(`
		Subject: noreply: Your password just changed.

		Hi {{.GetEmail}},
		Your password just change, this email notice you abou that.
		The new password is:
			{{.GetPassword}}
	`))
	n.infoChangedTmpl = template.Must(template.New("passChangedTmpl").Parse(`
		Subject: noreply: Your information just changed.

		Hi {{.GetEmail}},
		Your information just change, this email notice you abou that.
	`))
	n.privChangedTmpl = template.Must(template.New("passChangedTmpl").Parse(`
		Subject: noreply: Your privilege just changed.

		Hi {{.GetEmail}},
		Your privilege just change, this email notice you abou that.
	`))
	return n
}

func (n *SMTPNotificater) AccountAdded(user User) error {
	var buff bytes.Buffer
	n.accAddedTmpl.Execute(&buff, user)
	err := smtp.SendMail(n.Addr, n.Auth, n.Email, []string{user.GetEmail()}, buff.Bytes())
	return err
}

func (n *SMTPNotificater) PasswordChanged(user User) error {
	var buff bytes.Buffer
	n.passChangedTmpl.Execute(&buff, user)
	err := smtp.SendMail(n.Addr, n.Auth, n.Email, []string{user.GetEmail()}, buff.Bytes())
	return err
}

func (n *SMTPNotificater) AccountInfoChanged(user User) error {
	var buff bytes.Buffer
	n.infoChangedTmpl.Execute(&buff, user)
	err := smtp.SendMail(n.Addr, n.Auth, n.Email, []string{user.GetEmail()}, buff.Bytes())
	return err
}
func (n *SMTPNotificater) AccountPrivilegeChanged(user User) error {
	var buff bytes.Buffer
	n.privChangedTmpl.Execute(&buff, user)
	err := smtp.SendMail(n.Addr, n.Auth, n.Email, []string{user.GetEmail()}, buff.Bytes())
	return err
}

var _ Notificater = &SMTPNotificater{}

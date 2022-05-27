package mail

import (
	"context"
	"crypto/tls"
	"net"
	"net/smtp"
	"net/textproto"

	"github.com/jordan-wright/email"
	"github.com/pkg/errors"
)

// Mail struct holds necessary data to send emails.
type Mail struct {
	senderAddress     string
	smtpHostAddr      string
	smtpAuth          smtp.Auth
	receiverAddresses []string
	tls               bool
}

// New returns a new instance of a Mail notification service.
func New(senderAddress, smtpHostAddress string) *Mail {
	return &Mail{
		senderAddress:     senderAddress,
		smtpHostAddr:      smtpHostAddress,
		receiverAddresses: []string{},
	}
}

// AuthenticateSMTP authenticates you to send emails via smtp.
// Example values: "", "test@gmail.com", "password123", "smtp.gmail.com"
// For more information about smtp authentication, see here:
//    -> https://pkg.go.dev/net/smtp#PlainAuth
func (m *Mail) AuthenticateSMTP(identity, userName, password, host string) {
	m.smtpAuth = smtp.PlainAuth(identity, userName, password, host)
}

// LoginAuth authenticates you to send emails via smtp.
// Example values: "test@gmail.com", "password123"
// For more information about smtp authentication, see here:
//    -> https://pkg.go.dev/net/smtp#Auth
func (m *Mail) LoginAuth(userName, password string) {
	m.smtpAuth = &loginAuth{
		username: userName,
		password: password,
	}
}

// AddReceivers takes email addresses and adds them to the internal address list. The Send method will send
// a given message to all those addresses.
func (m *Mail) AddReceivers(addresses ...string) {
	m.receiverAddresses = append(m.receiverAddresses, addresses...)
}

// EnableTLS use tls/ssl to connect
func (m *Mail) EnableTLS() {
	m.tls = true
}

// Send takes a message subject and a message body and sends them to all previously set chats. Message body supports
// html as markup language.
func (m Mail) Send(ctx context.Context, subject, message string) error {
	msg := &email.Email{
		To:      m.receiverAddresses,
		From:    m.senderAddress,
		Subject: subject,
		// Text:    []byte("Text Body is, of course, supported!"),
		HTML:    []byte(message),
		Headers: textproto.MIMEHeader{},
	}

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	default:
		if m.tls {
			var serverName string
			serverName, _, err = net.SplitHostPort(m.smtpHostAddr)
			if err != nil {
				return err
			}
			err = msg.SendWithTLS(m.smtpHostAddr, m.smtpAuth, &tls.Config{ServerName: serverName})
		} else {
			err = msg.Send(m.smtpHostAddr, m.smtpAuth)
		}

		if err != nil {
			err = errors.Wrap(err, "failed to send mail")
		}
	}

	return err
}

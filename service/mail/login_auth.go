package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

// loginAuth implement smtp auth interface
type loginAuth struct {
	username, password string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		challenge := strings.ToLower(string(fromServer))
		switch challenge {
		case "username:":
			return []byte(a.username), nil
		case "password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unexpected server challenge, %s", challenge)
		}
	}
	return nil, nil
}

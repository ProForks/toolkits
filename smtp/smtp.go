package smtp

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

type loginAuth struct {
	username, password string
	address            string
}

func LoginAuth(username, password, host string) *loginAuth {
	return &loginAuth{username, password, host}
}

func (a loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if server.Name != a.address {
		return "", nil, errors.New("wrong host name")
	}
	return "LOGIN", nil, nil
}

func (a loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		if bytes.EqualFold([]byte("username:"), fromServer) {
			return []byte(a.username), nil
		} else if bytes.EqualFold([]byte("password:"), fromServer) {
			return []byte(a.password), nil
		}
	}
	return nil, nil
}

func (a *loginAuth) SendMail(from, tos, subject, body string, contentType ...string) error {
	if a.address == "" {
		return fmt.Errorf("address is necessary")
	}

	hp := strings.Split(a.address, ":")
	if len(hp) != 2 {
		return fmt.Errorf("address format error")
	}

	arr := strings.Split(tos, ";")
	count := len(arr)
	safeArr := make([]string, 0, count)
	for i := 0; i < count; i++ {
		if arr[i] == "" {
			continue
		}
		safeArr = append(safeArr, arr[i])
	}

	if len(safeArr) == 0 {
		return fmt.Errorf("tos invalid")
	}

	tos = strings.Join(safeArr, ";")

	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	header := make(map[string]string)
	header["From"] = from
	header["To"] = tos
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"

	ct := "text/plain; charset=UTF-8"
	if len(contentType) > 0 && contentType[0] == "html" {
		ct = "text/html; charset=UTF-8"
	}

	header["Content-Type"] = ct
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + b64.EncodeToString([]byte(body))
	auth := LoginAuth(a.username, a.password, hp[0])
	return smtp.SendMail(a.address, auth, from, strings.Split(tos, ";"), []byte(message))
}

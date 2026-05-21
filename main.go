package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/emersion/go-smtp"
)

type Backend struct{}

type Session struct {
	from string
	to   []string
}

func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	log.Printf("Connection from %s", c.Hostname())
	return &Session{}, nil
}

func (s *Session) AuthPlain(username, password string) error {
	if username != "admin" || password != "changeme" {
		return fmt.Errorf("invalid credentials")
	}
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.to = append(s.to, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	err = os.MkdirAll("/var/mailserver/mails", 0755)
	if err != nil {
		return err
	}

	filename := filepath.Join(
		"/var/mailserver/mails",
		fmt.Sprintf("%d.eml", time.Now().UnixNano()),
	)

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	log.Printf("Stored mail from %s to %v", s.from, s.to)

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func main() {

	be := &Backend{}

	s := smtp.NewServer(be)

	s.Addr = ":25"
	s.Domain = "mail.example.com"

	s.AllowInsecureAuth = false

	cert, err := tls.LoadX509KeyPair(
		"/usr/local/etc/mailserver/cert.pem",
		"/usr/local/etc/mailserver/key.pem",
	)

	if err != nil {
		log.Fatal(err)
	}

	s.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	log.Println("Mail server listening on port 25")

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

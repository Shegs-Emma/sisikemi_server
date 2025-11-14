package mail

import (
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain string
	Host string
	Port int
	Username string
	Password string
	Encryption string
	FromAddress string
	FromName string
}

type Message struct {
	Subject string
	VerifyUrl string
	Content string
	To string
	From string
	FromName string
}

type EmailSendSMTPEmailSender interface {
	SendSMTPEmail(message Message) error
}

func NewSendSMTPEmailSender(m Mail) *Mail {
	return &Mail{
		Domain: m.Domain,
		Host: m.Host,
		Port: m.Port,
		Username: m.Username,
		Password: m.Password,
		Encryption: m.Encryption,
		FromAddress: m.FromAddress,
		FromName: m.FromName,
	}
}

func (m *Mail) SendSMTPEmail(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)
	// email.SetBody(mail.TextPlain, msg.VerifyUrl)
	email.SetBody(mail.TextHTML, msg.Content)
	email.AddAlternative(mail.TextPlain, "Verify your email: "+msg.VerifyUrl)


	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
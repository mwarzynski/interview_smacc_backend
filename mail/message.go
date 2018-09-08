package mail

import (
	"github.com/badoux/checkmail"
	"github.com/pkg/errors"
)

// Message defines the message to be sent.
type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

// Validate checks if the message is valid.
func (m *Message) Validate() error {
	if err := checkmail.ValidateFormat(m.To); err != nil {
		return errors.Wrap(err, "mail 'to' is invalid")
	}
	if err := checkmail.ValidateFormat(m.From); err != nil {
		return errors.Wrap(err, "mail 'to' is invalid")
	}
	if m.Subject == "" {
		return errors.New("'subject' must not be empty")
	}
	if m.Text == "" {
		return errors.New("'text' must not be empty")
	}
	return nil
}

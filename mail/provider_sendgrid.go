package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type reqBody struct {
	Subject string    `json:"subject"`
	Content []content `json:"content"`

	From             email             `json:"from"`
	Personalizations []personalization `json:"personalizations"`
}

type content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type personalization struct {
	To []email `json:"to"`
}

type email struct {
	Email string `json:"email"`
}

type sendGridProvider struct {
	host   string
	apiKey string

	doer HTTPDoer
}

// NewSendGridProvider creates mail provided which uses SendGrid API to send mails.
func NewSendGridProvider(host, apiKey string, doer HTTPDoer) Provider {
	return &sendGridProvider{
		host:   host,
		doer:   doer,
		apiKey: apiKey,
	}
}

func (p *sendGridProvider) Name() string {
	return "sendgrid"
}

func (p *sendGridProvider) Send(ctx context.Context, message Message) error {
	sendURL := url.URL{
		Scheme: "https",
		Host:   p.host,
		Path:   "/v3/mail/send",
	}

	data := reqBody{
		Subject: message.Subject,
		Content: []content{
			content{
				Type:  "text/plain",
				Value: message.Text,
			},
		},
		From: email{
			Email: message.From,
		},
		Personalizations: []personalization{
			personalization{
				To: []email{
					email{
						Email: message.To,
					},
				},
			},
		},
	}

	rawData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "couldn't marshal body")
	}

	req, err := http.NewRequest(http.MethodPost, sendURL.String(), bytes.NewReader(rawData))
	if err != nil {
		return errors.Wrap(err, "couldn't create new request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))
	req.WithContext(ctx)

	resp, err := p.doer.Do(req)
	if err != nil {
		return errors.Wrap(err, "couldn't do request")
	}
	if resp == nil { // this case is possible when system network is down
		return errors.New("response is nil")
	}

	if resp.StatusCode != http.StatusAccepted {
		var body []byte
		if resp.Body != nil {
			body, _ = ioutil.ReadAll(io.LimitReader(resp.Body, 1024*20)) // 20 kb
		}
		return errors.Errorf("invalid status code %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

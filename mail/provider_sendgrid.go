package mail

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type sendGridProvider struct {
	host    string
	apiUser string
	apiKey  string

	doer HTTPDoer
}

// NewSendGridProvider creates mail provided which uses SendGrid API to send mails.
func NewSendGridProvider(host, apiKey, apiUser string, doer HTTPDoer) Provider {
	return &sendGridProvider{
		host:    host,
		doer:    doer,
		apiKey:  apiKey,
		apiUser: apiUser,
	}
}

func (p *sendGridProvider) Name() string {
	return "sendgrid"
}

func (p *sendGridProvider) Send(ctx context.Context, message Message) error {
	sendURL := url.URL{
		Scheme:  "https",
		Host:    p.host,
		RawPath: "/api/mail.send.json",
	}

	form := url.Values{
		"api_user": []string{p.apiUser},
		"api_key":  []string{p.apiKey},
		"from":     []string{message.From},
		"to":       []string{message.To},
		"subject":  []string{message.Subject},
		"text":     []string{message.Text},
	}

	req, err := http.NewRequest(http.MethodPost, sendURL.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return errors.Wrap(err, "couldn't create new request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.WithContext(ctx)

	resp, err := p.doer.Do(req)
	if err != nil {
		return errors.Wrap(err, "couldn't do request")
	}
	if resp == nil { // this case is possible when system network is down
		return errors.New("response is nil")
	}

	if resp.StatusCode != http.StatusOK {
		var body []byte
		if resp.Body != nil {
			body, _ = ioutil.ReadAll(io.LimitReader(resp.Body, 1024*20)) // 20 kb
		}
		return errors.Errorf("invalid status code %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

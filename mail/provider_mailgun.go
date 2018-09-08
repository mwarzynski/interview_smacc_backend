package mail

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type mailGunProvider struct {
	host   string
	apiKey string
	domain string

	doer HTTPDoer
}

// NewMailGunProvider creates mail provided which uses MailGun API to send mails.
// Visit https://mailgun.net for more information.
func NewMailGunProvider(host, apiKey, domain string, doer HTTPDoer) Provider {
	return &mailGunProvider{
		host:   host,
		apiKey: apiKey,
		domain: domain,
		doer:   doer,
	}
}

func (p *mailGunProvider) Name() string {
	return "mailgun"
}

func (p *mailGunProvider) Send(ctx context.Context, message Message) error {
	sendURL := url.URL{
		Scheme:  "https",
		Host:    p.host,
		RawPath: fmt.Sprintf("/v3/%s/messages", p.domain),
	}
	form := url.Values{
		"from":    []string{message.From},
		"to":      []string{message.To},
		"subject": []string{message.Subject},
		"text":    []string{message.Text},
	}

	req, err := http.NewRequest(http.MethodPost, sendURL.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return errors.Wrap(err, "couldn't create new request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(
		"Authorization",
		fmt.Sprintf(
			"Basic %s",
			base64.StdEncoding.EncodeToString([]byte("api:"+p.apiKey)),
		),
	)
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

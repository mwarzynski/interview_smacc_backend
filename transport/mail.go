package transport

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/mwarzynski/smacc-backend/mail"
)

// Mail is a handler for sending mail messages.
func Mail(mailService *mail.Service, l logrus.FieldLogger) http.HandlerFunc {
	l = l.WithField("tags", []string{"handler", "mail"})

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024)) // 1 MB
		if err != nil {
			if err == io.EOF {
				http.Error(w, "body size must be <= 1 MB", http.StatusBadRequest)
				return
			}
			err = errors.Wrap(err, "couldn't read the body")
			l.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var args mail.SendArgs
		if err := json.Unmarshal(body, &args); err != nil {
			err = errors.Wrap(err, "body must be a valid JSON")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sent := mailService.Send(r.Context(), args)
		if !sent {
			http.Error(w, "couldn't send the mail", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

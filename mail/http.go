package mail

import "net/http"

// HTTPDoer handles http request.
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

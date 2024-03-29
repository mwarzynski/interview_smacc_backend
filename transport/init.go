package transport

import (
	"context"
	"net/http"

	gometrics "github.com/armon/go-metrics"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/mwarzynski/smacc-backend/mail"
	"github.com/sirupsen/logrus"
)

// HTTPDoer is an object that makes HTTP requests (*http.Client).
type HTTPDoer interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

// Init constructs http.Handler for the main path.
func Init(
	appCtx context.Context,
	mailService *mail.Service,
	metricsService *gometrics.Metrics,
	l logrus.FieldLogger,
) http.Handler {
	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.With(metricsMiddleware("mail-send", metricsService)).Post("/mail/send", Mail(mailService, l))
	})

	r.Get("/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	return r
}

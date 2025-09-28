package app

import (
	"net/http"

	memes "github.com/ecerizola-im/AnnoyEm/internal/memes"
)

func Router(svc *memes.MemeService) http.Handler {
	mux := http.NewServeMux()

	h := memes.NewHandler(svc)
	h.Register(mux)

	return LoggingMiddleware(mux)
}

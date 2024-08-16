package middleware

import (
	"compress/gzip"
	"net/http"
	"runtime/debug"

	sugared "github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

func LogMiddleware(next http.HandlerFunc, logger sugared.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infow(
			"entering",
			"method", r.Method,
			"url", r.RequestURI,
		)

		next(w, r)

		logger.Infow(
			"leaving",
			"method", r.Method,
			"url", r.RequestURI,
		)
	}
}

func PanicMiddleware(next http.HandlerFunc, logger sugared.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorw(
					"panic happened",
					"reason", err,
					"stacktrace", string(debug.Stack()),
				)

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next(w, r)
	}
}

func DecompressMiddleware(next http.HandlerFunc, logger sugared.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request body", http.StatusBadRequest)
				return
			}
			defer gz.Close()

			r.Body = gz
		}

		next(w, r)
	}
}

func ErrorMiddleware(next func(http.ResponseWriter, *http.Request) error, logger sugared.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			logger.Errorw(
				"error happened",
				"url", r.RequestURI,
				"reason", err,
			)

			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func Accept(f func(http.ResponseWriter, *http.Request) error, logger sugared.Logger) http.HandlerFunc {
	middlewares := []func(next http.HandlerFunc, logger sugared.Logger) http.HandlerFunc{
		PanicMiddleware,
		LogMiddleware,
		DecompressMiddleware,
	}

	prelude := ErrorMiddleware(f, logger)

	for i := range middlewares {
		prelude = middlewares[i](prelude, logger)
	}

	return prelude
}

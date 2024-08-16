package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"runtime/debug"
	"strings"

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

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CompressMiddleware(next http.HandlerFunc, logger sugared.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		next(gzipWriter{ResponseWriter: w,
			Writer: gz,
		}, r)
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
		CompressMiddleware,
	}

	prelude := ErrorMiddleware(f, logger)

	for i := range middlewares {
		prelude = middlewares[i](prelude, logger)
	}

	return prelude
}

package impl

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
	"github.com/test-go/testify/assert"

	businessimpl "github.com/chernyshevuser/practicum-metrics-collector/internal/server/business/impl"
	storageimpl "github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage/impl"
)

func TestUpdateMetric(t *testing.T) {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	logger := logger.New()
	defer logger.Sync()

	dbSvc, err := storageimpl.New(mainCtx, logger)
	if err != nil {
		logger.Errorw(
			"cant create db svc",
			"reason", err,
		)
		panic("db initialization failed")
	}
	defer dbSvc.Close()

	businessSvc := businessimpl.New(dbSvc, logger)
	defer businessSvc.Close()

	svc := New(businessSvc, logger)

	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		endp string
		want want
	}{
		{
			name: "correct counter",
			endp: "/update/counter/sample_text/123",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "incorrect counter by value",
			endp: "/update/counter/sample_text/123.5",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "correct gauge with int value",
			endp: "/update/gauge/sample_text/123",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "correct gauge with float value",
			endp: "/update/gauge/name/123.4",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "incorrect gauge by value",
			endp: "/update/gauge/sample_text/value",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.endp, nil)
			w := httptest.NewRecorder()
			svc.UpdateMetric(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

package impl

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	mockbusiness "github.com/chernyshevuser/practicum-metrics-collector/internal/server/business/mock"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/router"
	mocklogger "github.com/chernyshevuser/practicum-metrics-collector/tools/logger/mock"
	"github.com/golang/mock/gomock"
)

// Example_api_PingDB demonstrates how to use PingDB handler.
func Example_api_PingDB() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	businessSvc := mockbusiness.NewMockMetricsCollector(ctrl)
	businessSvc.EXPECT().PingDB(gomock.Any()).Return(nil)

	svc := New(businessSvc, logger)

	req := httptest.NewRequest(http.MethodGet, router.PingDB, nil)
	w := httptest.NewRecorder()

	svc.PingDB(w, req)

	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	// Output:
	// 200
}

package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	businessimpl "github.com/chernyshevuser/practicum-metrics-collector/internal/server/business/impl"
	api "github.com/chernyshevuser/practicum-metrics-collector/internal/server/handler/impl"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/router"
	storageimpl "github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage/impl"
	"go.uber.org/zap"

	_ "net/http/pprof"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/config"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
)

func main() {
	printVersion()

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	logger := zap.Must(zap.NewProductionConfig().Build()).Sugar()
	defer logger.Sync()

	config.Setup(logger)

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

	apiSvc := api.New(businessSvc, logger)

	muxRouter := mux.NewRouter()
	router.SetupRouter(apiSvc, muxRouter, logger, config.CryptoKey)
	muxRouter.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	server := http.Server{
		Addr:    config.Addr,
		Handler: muxRouter,
		BaseContext: func(_ net.Listener) context.Context {
			return mainCtx
		},
	}

	g, gCtx := errgroup.WithContext(mainCtx)
	g.Go(func() error {
		return server.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return server.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		logger.Infow(
			"server exit",
			"reason", err,
		)
	}
}

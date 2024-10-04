package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business/impl"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/config"
	"go.uber.org/zap"
)

func main() {
	printVersion()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	logger := zap.Must(zap.NewProductionConfig().Build()).Sugar()
	defer logger.Sync()

	config.Setup(logger)

	agentSvc := impl.New(
		logger,
		config.PollInterval,
		config.ReportInterval,
		config.HashKey,
		fmt.Sprintf("http://%s/updates/", config.Addr),
		config.RateLimit,
	)
	defer agentSvc.Close()

	agentSvc.Run(ctx)
}

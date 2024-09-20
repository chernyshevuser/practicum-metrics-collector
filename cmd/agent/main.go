package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business/impl"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/config"
	logger "github.com/chernyshevuser/practicum-metrics-collector/tools/logger/impl"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	logger := logger.New()
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

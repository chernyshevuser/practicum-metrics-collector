package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business/impl"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/config"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	logger := logger.New()
	defer logger.Sync()

	//TODO fix
	config.Setup(logger)

	agentSvc := impl.New(
		logger,
		config.PollInterval,
		config.ReportInterval,
		config.HashKey,
		fmt.Sprintf("http://%s/update/", config.Addr),
	)
	defer agentSvc.Close()

	agentSvc.Run(ctx)
}

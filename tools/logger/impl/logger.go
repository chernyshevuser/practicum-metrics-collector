package impl

import (
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(opts ...zap.Option) logger.Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = encoderConfig()

	return zap.Must(cfg.Build(withStdOpts(opts)...)).Sugar()
}

func withStdOpts(opts []zap.Option) []zap.Option {
	return append(opts, zap.AddStacktrace(zap.PanicLevel))
}

func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochMillisTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

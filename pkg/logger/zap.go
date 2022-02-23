package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newZapLogger(config *Config) (*Logger, error) {
	l := new(Logger)

	cores := newCores(l, config)
	combinedCore := zapcore.NewTee(cores...)
	options := getZapOptions(config)

	l.raw = zap.New(combinedCore, options...)
	l.sugared = l.raw.WithOptions(zap.AddCallerSkip(1)).Sugar()

	// defers.Register(func() {
	// 	_ = l.sugared.Sync()
	// 	_ = l.raw.Sync()
	// })

	return l, nil
}

func newCores(l *Logger, config *Config) []zapcore.Core {
	var cores []zapcore.Core

	if config.EnableConsole {
		ws := newConsoleWriterSyncer(config)
		l.consoleLevel = zap.NewAtomicLevelAt(getZapLevel(config.ConsoleLevel))
		core := zapcore.NewCore(getEncoder(config.ConsoleJSONFormat), ws, l.consoleLevel)
		cores = append(cores, core)
	}

	if config.EnableFile {
		ws := newFileWriteSyncer(config)
		l.fileLevel = zap.NewAtomicLevelAt(getZapLevel(config.FileLevel))
		core := zapcore.NewCore(getEncoder(config.FileJSONFormat), ws, l.fileLevel)
		cores = append(cores, core)
	}

	return cores
}

func newConsoleWriterSyncer(config *Config) zapcore.WriteSyncer {
	if config.ConsoleWriteAsync {
		return newWriteAsyncer(os.Stdout)
	}

	return zapcore.Lock(os.Stdout)
}

func newFileWriteSyncer(config *Config) zapcore.WriteSyncer {
	writer := &lumberjack.Logger{
		Filename:   config.FileName,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxFileHistory,
		Compress:   config.Compress,
	}

	if config.FileWriteAsync {
		return newWriteAsyncer(writer)
	}

	return zapcore.AddSync(writer)
}

func getZapOptions(config *Config) []zap.Option {
	var options []zap.Option
	options = append(options, zap.AddStacktrace(zapcore.DPanicLevel))

	if config.AddCaller {
		options = append(options, zap.AddCaller())
	}

	// if env.DeployEnv == env.DeployEnvDev {
	// 	options = append(options, zap.Development())
	// }

	return options
}

func getEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "log.level",
		NameKey:        "logger",
		CallerKey:      "log.origin.file.name",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

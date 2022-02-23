package logger

import (
	"errors"

	"go.uber.org/zap"
)

type Config struct {
	LoggerType        string `default:"zap" env:"LOG_LOGGER_TYPE"`
	EnableConsole     bool   `default:"true" env:"LOG_ENABLE_STDOUT"`
	ConsoleWriteAsync bool   `default:"false" env:"LOG_CONSOLE_WRITE_ASYNC"`
	ConsoleJSONFormat bool   `default:"true" env:"LOG_CONSOLE_JSON_FORMAT"`
	ConsoleLevel      string `default:"debug" env:"LOG_CONSOLE_LEVEL"`
	EnableFile        bool   `default:"true" env:"LOG_ENABLE_FILE"`
	FileWriteAsync    bool   `default:"true" env:"LOG_FILE_WRITE_ASYNC"`
	FileJSONFormat    bool   `default:"false" env:"LOG_FILE_JSON_FORMAT"`
	FileLevel         string `default:"debug" env:"LOG_FILE_LEVEL"`
	FileName          string `default:"unknown" env:"LOG_FILE_NAME"`
	MaxSize           int    `default:"20" env:"LOG_MAX_SIZE"`
	MaxBackups        int    `default:"100" env:"LOG_MAX_BACKUPS"`
	Compress          bool   `default:"true" env:"LOG_ENABLE_COMPRESS"`
	AddCaller         bool   `default:"true" env:"LOG_ADD_CALLER"`
	MaxFileHistory    int    `default:"30" env:"LOG_MAX_FILE_HISTORY"`
}

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	FatalLevel = "fatal"
)

type Logger struct {
	raw          *zap.Logger
	sugared      *zap.SugaredLogger
	consoleLevel zap.AtomicLevel
	fileLevel    zap.AtomicLevel
}

type Fields map[string]interface{}

var (
	_defaultConfig = Config{
		LoggerType:        "zap",
		EnableConsole:     true,
		ConsoleWriteAsync: true,
		ConsoleJSONFormat: false,
		ConsoleLevel:      "debug",
		EnableFile:        true,
		FileWriteAsync:    true,
		FileJSONFormat:    false,
		FileLevel:         "debug",
		FileName:          "./log/notionSync.log",
		MaxSize:           20,
		MaxBackups:        100,
		Compress:          false,
		AddCaller:         true,
		MaxFileHistory:    30,
	}
)

var _log *Logger

func Debugf(format string, args ...interface{}) {
	_log.sugared.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	_log.sugared.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	_log.sugared.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	_log.sugared.Errorf(format, args...)
}

func DPanicf(format string, args ...interface{}) {
	_log.sugared.DPanicf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	_log.sugared.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	_log.sugared.Fatalf(format, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	_log.sugared.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	_log.sugared.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	_log.sugared.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	_log.sugared.Errorw(msg, keysAndValues...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	_log.sugared.DPanicw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	_log.sugared.Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	_log.sugared.Fatalw(msg, keysAndValues...)
}

func WithFields(fields Fields) *Logger {
	var f = make([]interface{}, 0, len(fields))
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}

	sugar := _log.raw.WithOptions(zap.AddCallerSkip(1)).Sugar().With(f...)
	return &Logger{
		sugared: sugar,
	}
}

func Zap() *zap.Logger {
	return _log.raw
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sugared.Debugf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.sugared.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.sugared.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.sugared.Errorf(format, args...)
}

func (l *Logger) DPanicf(format string, args ...interface{}) {
	l.sugared.DPanicf(format, args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.sugared.Panicf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.sugared.Fatalf(format, args...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.sugared.Debugw(msg, keysAndValues...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugared.Infow(msg, keysAndValues...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.sugared.Warnw(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugared.Errorw(msg, keysAndValues...)
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.sugared.DPanicw(msg, keysAndValues...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.sugared.Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.sugared.Fatalw(msg, keysAndValues...)
}

func (l *Logger) Zap() *zap.Logger {
	return l.raw
}

func SetEnableLevel(outPut, level string) (ok bool) {
	return _log.SetEnableLevel(outPut, level)
}

func (l *Logger) SetEnableLevel(outPut, level string) (ok bool) {
	switch outPut {
	case "console":
		l.consoleLevel.SetLevel(getZapLevel(level))
		return true
	case "file":
		l.fileLevel.SetLevel(getZapLevel(level))
		return true
	}
	return true
}

func Init(config *Config) {
	switch config.LoggerType {
	case "zap":
		logger, err := newZapLogger(config)
		if err != nil {
			panic(err)
		}
		_log = logger
	default:
		panic(errors.New("invalid Logger type"))
	}
}

func NewLogger(config *Config) *Logger {
	l, err := newZapLogger(config)
	if err != nil {
		panic(err)
	}
	return l
}

func init() {
	logger, err := newZapLogger(&_defaultConfig)
	if err != nil {
		panic(err)
	}
	_log = logger
}

package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logCtxKey int

const _loggerKey logCtxKey = iota

func WithField(key, value string) zapcore.Field {
	return zap.String(key, value)
}

func WithTrace(ctx context.Context, fields ...zapcore.Field) context.Context {
	return context.WithValue(ctx, _loggerKey, T(ctx).with(fields...))
}

func SetTraceFromOldContext(ctx context.Context, oldCtx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if oldCtx == nil {
		return ctx
	}

	return context.WithValue(ctx, _loggerKey, T(oldCtx))
}

func T(ctx context.Context) *Logger {
	if ctx == nil {
		goto end
	}

	if ctxLogger, ok := ctx.Value(_loggerKey).(*Logger); ok {
		return ctxLogger
	}

end:
	return _log
}

func (l *Logger) with(keyValues ...zapcore.Field) *Logger {
	var args = make([]interface{}, 0, len(keyValues))
	for _, keyValue := range keyValues {
		args = append(args, keyValue)
	}

	rv := new(Logger)
	*rv = *l
	rv.sugared = l.raw.WithOptions(zap.AddCallerSkip(1)).Sugar().With(args...)

	return rv
}

package logger

import "context"

type logKey struct{}

func InjectLogToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logKey{}, Call())
}

func LoadLogFromContext(ctx context.Context) *LogPayload {
	getLog := ctx.Value(logKey{})

	if log, ok := getLog.(*LogPayload); ok {
		return log
	}

	return Call()
}

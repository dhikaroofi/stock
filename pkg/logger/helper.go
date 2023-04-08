package logger

import "context"

func ExtractLogFromContext(ctx context.Context) *LogPayload {
	getLog := ctx.Value(LogKey)

	if log, ok := getLog.(*LogPayload); ok {
		return log
	}

	return Call()
}

package log

import "context"

var gLogger *Logger

func SetupGlobal(c *Config) error {
	l, err := NewLogger(c)
	if err != nil {
		return err
	}
	if gLogger != nil {
		if err := gLogger.Close(); err != nil {
			l.Warn(context.Background(), "close global logger failed", "err", err)
		}
	}
	gLogger = l
	return nil
}

func Debug(ctx context.Context, msg string, fields ...interface{}) {
	gLogger.Debug(ctx, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...interface{}) {
	gLogger.Info(ctx, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...interface{}) {
	gLogger.Warn(ctx, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...interface{}) {
	gLogger.Error(ctx, msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...interface{}) {
	gLogger.Fatal(ctx, msg, fields...)
}

func Panic(ctx context.Context, msg string, fields ...interface{}) {
	gLogger.Fatal(ctx, msg, fields...)
}

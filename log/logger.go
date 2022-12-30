package log

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	// lumberjack logger config
	Filename   string `yaml:"Filename"`
	MaxSize    int    `yaml:"Maxsize"`
	MaxAge     int    `yaml:"Maxage"`
	MaxBackups int    `yaml:"Maxbackups"`
	LocalTime  bool   `yaml:"Localtime"`
	Compress   bool   `yaml:"Compress"`

	CallerSkip int    `yaml:"CallerSkip"`
	Level      string `yaml:"Level"`
	Console    string `yaml:"Console"`

	// gorm log config
	GormLevel                 string        `yaml:"GormLevel"`
	SqlSlowThreshold          time.Duration `yaml:"SqlSlowThreshold"`
	IgnoreRecordNotFoundError bool          `yaml:"IgnoreRecordNotFoundError"`
	IgnoreDuplicateError      bool          `yaml:"IgnoreDuplicateError"`
}

type HookFunc func(ctx context.Context, level zapcore.Level, msg string, fields []zap.Field)

type Logger struct {
	c     *Config
	z     *zap.Logger
	l     zapcore.Level
	hooks []HookFunc
}

func newZapCore(c *Config, level zapcore.Level) (zapcore.Core, error) {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:       "M",
		LevelKey:         "L",
		TimeKey:          "T",
		NameKey:          "N",
		CallerKey:        "C",
		FunctionKey:      "F",
		StacktraceKey:    "S",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		EncodeDuration:   zapcore.StringDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: "\t",
	}

	syncers := make([]zapcore.WriteSyncer, 0, 2)
	if len(c.Filename) > 0 {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   c.Filename,
			MaxSize:    c.MaxSize,
			MaxAge:     c.MaxAge,
			MaxBackups: c.MaxBackups,
			LocalTime:  c.LocalTime,
			Compress:   c.Compress,
		}
		syncers = append(syncers, zapcore.AddSync(lumberjackLogger))
	}
	switch strings.ToLower(c.Console) {
	case "stdout", "1":
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	case "stderr", "2":
		syncers = append(syncers, zapcore.AddSync(os.Stderr))
	}
	if len(syncers) == 0 {
		return nil, fmt.Errorf("no filename or valid console in config")
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.NewMultiWriteSyncer(syncers...), level)
	return core, nil
}

type writeHook struct{}

func (w writeHook) OnWrite(*zapcore.CheckedEntry, []zapcore.Field) {
}

func NewLogger(c *Config, cores ...zapcore.Core) (*Logger, error) {
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(c.Level)); err != nil {
		return nil, err
	}
	if len(cores) == 0 {
		core, err := newZapCore(c, level)
		if err != nil {
			return nil, err
		}
		cores = []zapcore.Core{core}
	}
	options := []zap.Option{
		zap.WithCaller(true),
		zap.AddCallerSkip(c.CallerSkip),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.WithFatalHook(writeHook{}),
	}

	cfg := *c
	return &Logger{
		c: &cfg,
		z: zap.New(zapcore.NewTee(cores...), options...),
		l: level,
	}, nil
}

func appendFields(fields []zap.Field, kvs ...interface{}) []zap.Field {
	if len(kvs) == 1 {
		if arr, ok := kvs[0].([]interface{}); ok {
			kvs = arr
		}
	}
	newFields := make([]zap.Field, len(fields), len(fields)+(len(kvs)+1)/2)
	copy(newFields, fields)

	for i := 0; i < len(kvs)-1; i += 2 {
		k, v := kvs[i], kvs[i+1]
		if str, ok := k.(string); ok {
			newFields = append(newFields, zap.Any(str, v))
		} else {
			newFields = append(newFields, zap.Any("invalidKey", v))
		}
	}
	if len(kvs)%2 != 0 { //the last value without a key
		newFields = append(newFields, zap.Any("danglingKey", kvs[len(kvs)-1]))
	}
	return newFields
}

func (l *Logger) Core() zapcore.Core {
	return l.z.Core()
}

// NewLevel gorm log need
func (l *Logger) NewLevel(lv zapcore.Level) *Logger {
	return &Logger{
		z:     l.z,
		l:     lv,
		hooks: l.hooks,
	}
}

func (l *Logger) AddHook(hook HookFunc) {
	l.hooks = append(l.hooks, hook)
}

type contextLogKey struct{}

func (l *Logger) WithValues(ctx context.Context, kvs ...interface{}) context.Context {
	logFields, _ := ctx.Value(contextLogKey{}).([]zap.Field)
	return context.WithValue(ctx, contextLogKey{}, appendFields(logFields, kvs...))
}

func (l *Logger) Log(ctx context.Context, level zapcore.Level, msg string, kvs ...interface{}) {
	if !l.l.Enabled(level) {
		return
	}

	logFields, _ := ctx.Value(contextLogKey{}).([]zap.Field)
	fields := appendFields(logFields, kvs)
	switch level {
	case zapcore.DebugLevel:
		l.z.Debug(msg, fields...)
	case zapcore.InfoLevel:
		l.z.Info(msg, fields...)
	case zapcore.WarnLevel:
		l.z.Warn(msg, fields...)
	case zapcore.ErrorLevel:
		l.z.Error(msg, fields...)
	case zapcore.FatalLevel:
		l.z.Fatal(msg, fields...)
	case zap.PanicLevel:
		l.z.Fatal(msg, fields...)
	}
	for i := range l.hooks {
		l.hooks[i](ctx, level, msg, fields)
	}
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	l.Log(ctx, zapcore.DebugLevel, msg, fields...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...interface{}) {
	l.Log(ctx, zapcore.InfoLevel, msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	l.Log(ctx, zapcore.WarnLevel, msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...interface{}) {
	l.Log(ctx, zapcore.ErrorLevel, msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...interface{}) {
	l.Log(ctx, zapcore.FatalLevel, msg, fields...)
}

func (l *Logger) Panic(ctx context.Context, msg string, fields ...interface{}) {
	l.Log(ctx, zapcore.PanicLevel, msg, fields...)
}

func (l *Logger) Close() error {
	return l.z.Sync()
}

package log

import (
	"context"
	"errors"
	"fmt"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

// depend on global logger

func strToGormLogLevel(str string) gormlog.LogLevel {
	switch str {
	case "error", "ERROR":
		return gormlog.Error
	case "info", "INFO":
		return gormlog.Info
	default:
		return gormlog.Warn
	}
}

func NewGormLogger(config *Config) (gormlog.Interface, error) {
	c := *config
	level := strToGormLogLevel(config.GormLevel)

	if c.CallerSkip <= 2 {
		c.CallerSkip = 3
	}
	l, err := NewLogger(&c, gLogger.Core())
	if err != nil {
		return nil, err
	}
	return &GormLogger{
		logger:                    l,
		level:                     level,
		slowThreshold:             c.SqlSlowThreshold,
		ignoreRecordNotFoundError: c.IgnoreRecordNotFoundError,
		ignoreDuplicateError:      c.IgnoreDuplicateError,
	}, nil

}

type GormLogger struct {
	logger                    *Logger
	level                     gormlog.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
	ignoreDuplicateError      bool
}

func (g *GormLogger) LogMode(level gormlog.LogLevel) gormlog.Interface {
	var l zapcore.Level
	switch level {
	case gormlog.Info:
		l = zapcore.InfoLevel
	case gormlog.Warn:
		l = zapcore.WarnLevel
	case gormlog.Error:
		l = zapcore.ErrorLevel
	case gormlog.Silent:
		l = zapcore.FatalLevel
	default:
		Warn(context.TODO(), "invalid gorm log level", "level", level)
		return g
	}
	return &GormLogger{
		logger:                    g.logger.NewLevel(l),
		level:                     level,
		slowThreshold:             g.slowThreshold,
		ignoreDuplicateError:      g.ignoreDuplicateError,
		ignoreRecordNotFoundError: g.ignoreRecordNotFoundError,
	}
}

func (g *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	g.logger.Info(ctx, s, i...)
}

func (g *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	g.logger.Warn(ctx, s, i...)
}

func (g *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	g.logger.Error(ctx, s, i...)
}

func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.level <= gormlog.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && g.level >= gormlog.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !g.ignoreRecordNotFoundError) && (!IsDuplicateEntry(err) || !g.ignoreDuplicateError):
		sql, rows := fc()
		if rows == -1 {
			g.Error(ctx, err.Error(), "cost", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			g.Error(ctx, err.Error(), "cost", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	case elapsed > g.slowThreshold && g.slowThreshold != 0 && g.level >= gormlog.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", g.slowThreshold)
		if rows == -1 {
			g.Warn(ctx, slowLog, "cost", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			g.Warn(ctx, slowLog, "cost", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	case g.level == gormlog.Info:
		sql, rows := fc()
		if rows == -1 {
			g.Info(ctx, "gorm", "cost", float64(elapsed.Nanoseconds())/1e6, "sql", sql)
		} else {
			g.Info(ctx, "gorm", "cost", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	}
}

func IsDuplicateEntry(err error) bool {
	if e, ok := err.(*mysqlDriver.MySQLError); ok {
		return e.Number == 1062
	}

	return false
}

package log

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	Logger   *zap.Logger
	LogLevel logger.LogLevel
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Logger.Sugar().Infof(s, i...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Logger.Sugar().Warnf(s, i...)
	}
}

func (l *GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Logger.Sugar().Errorf(s, i...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	if err != nil {
		l.Logger.Error("SQL 执行错误", append(fields, zap.Error(err))...)
	} else if elapsed > 200*time.Millisecond { // 慢查询阈值
		l.Logger.Warn("慢查询", fields...)
	} else {
		l.Logger.Debug("SQL 执行", fields...)
	}
}

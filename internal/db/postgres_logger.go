/*
 * MIT License
 *
 * Copyright (c) 2026 nuoxi4n, NuoXianTech
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package db

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/NuoXianTech/api/internal/logger"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type gormZapLogger struct {
	logLevel                  gormLogger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

func (l *gormZapLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	copy := *l
	copy.logLevel = level
	return &copy
}

func (l *gormZapLogger) Info(ctx context.Context, message string, args ...interface{}) {
	if l.logLevel >= gormLogger.Info {
		logger.InfoF(ctx, message, args...)
	}
}

func (l *gormZapLogger) Warn(ctx context.Context, message string, args ...interface{}) {
	if l.logLevel >= gormLogger.Warn {
		logger.WarnF(ctx, message, args...)
	}
}

func (l *gormZapLogger) Error(ctx context.Context, message string, args ...interface{}) {
	if l.logLevel >= gormLogger.Error {
		logger.ErrorF(ctx, message, args...)
	}
}

// ParamsFilter 保留 SQL 占位符，避免将绑定参数展开到日志。
func (l *gormZapLogger) ParamsFilter(_ context.Context, sql string, _ ...interface{}) (string, []interface{}) {
	return sql, nil
}

func (l *gormZapLogger) Trace(ctx context.Context, begin time.Time, query func() (string, int64), err error) {
	elapsed := time.Since(begin)
	failed := err != nil && !(l.ignoreRecordNotFoundError && errors.Is(err, gorm.ErrRecordNotFound))
	slow := l.slowThreshold > 0 && elapsed > l.slowThreshold
	var write func(context.Context, string, ...interface{})
	switch {
	case failed && l.logLevel >= gormLogger.Error:
		write = logger.ErrorF
	case slow && l.logLevel >= gormLogger.Warn:
		write = logger.WarnF
	case l.logLevel >= gormLogger.Info:
		write = logger.InfoF
	default:
		return
	}
	sql, rows := query()
	write(ctx, "[PostgreSQL] duration=%s rows=%d slow=%t error=%v sql=%s", elapsed, rows, slow, err, sql)
}

func parseLogLevel(level string) gormLogger.LogLevel {
	switch strings.ToLower(level) {
	case "silent":
		return gormLogger.Silent
	case "error":
		return gormLogger.Error
	case "warn":
		return gormLogger.Warn
	default:
		return gormLogger.Info
	}
}

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

package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/NuoXianTech/api/internal/config"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logWriter         zapcore.WriteSyncer
	initLogWriterOnce sync.Once
	initLogWriterErr  error
)

// GetLogWriter 获取日志输出写入器
func GetLogWriter() (zapcore.WriteSyncer, error) {
	initLogWriterOnce.Do(func() {
		logWriter, initLogWriterErr = initWriter()
	})

	return logWriter, initLogWriterErr
}

func initWriter() (zapcore.WriteSyncer, error) {
	logConfig := config.Config.Log

	switch logConfig.Output {
	case "stdout":
		return zapcore.AddSync(os.Stdout), nil
	case "file":

		// 初始化日志目录
		logPath := logConfig.FilePath
		if strings.TrimSpace(logPath) == "" {
			return nil, fmt.Errorf("[Logger] log.file_path must not be empty")
		}
		logDir := filepath.Dir(logPath)
		if err := os.MkdirAll(logDir, 0750); err != nil {
			return nil, fmt.Errorf("[Logger] create log file dir err: %w", err)
		}

		// 配置日志轮转
		logOutput := &lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    logConfig.MaxSize,
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
		}

		return zapcore.AddSync(logOutput), nil
	default:
		return nil, fmt.Errorf("[Logger] invalid log output: %s", logConfig.Output)
	}
}

// getEncoder 获取日志编码器
func getEncoder() (zapcore.Encoder, error) {
	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	switch config.Config.Log.Format {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig), nil
	case "text":
		return zapcore.NewConsoleEncoder(encoderConfig), nil
	default:
		return nil, fmt.Errorf("[Logger] invalid log format: %s", config.Config.Log.Format)
	}
}

// getLogLevel 获取日志级别
func getLogLevel() (zapcore.Level, error) {
	level := config.Config.Log.Level

	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	case "panic":
		return zapcore.PanicLevel, nil
	case "fatal":
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("[Logger] invalid log level: %s", level)
	}
}

func getTraceIDFields(ctx context.Context) []zap.Field {
	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()
	return []zap.Field{
		zap.String("traceID", spanContext.TraceID().String()),
		zap.String("spanID", spanContext.SpanID().String()),
	}
}

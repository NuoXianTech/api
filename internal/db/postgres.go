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
	"log"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/NuoXianTech/api/internal/config"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	db *gorm.DB
)

func init() {
	if !config.Config.Database.Enabled {
		log.Println("[PostgreSQL] is disabled, skipping PostgreSQL initialization")
		return
	}

	dbConfig := config.Config.Database
	if dbConfig.Host == "" || dbConfig.Username == "" || dbConfig.Database == "" || dbConfig.Port < 1 || dbConfig.Port > 65535 {
		log.Fatal("[PostgreSQL] invalid host, username, database or port")
	}
	if dbConfig.MaxIdleConn < 0 || dbConfig.MaxOpenConn < 0 || dbConfig.ConnMaxLifetime < 0 {
		log.Fatal("[PostgreSQL] connection pool settings must not be negative")
	}

	// URL 编码保留账号、密码和数据库名中的特殊字符。
	dsn := url.URL{
		Scheme:  "postgres",
		User:    url.UserPassword(dbConfig.Username, dbConfig.Password),
		Host:    net.JoinHostPort(dbConfig.Host, strconv.Itoa(dbConfig.Port)),
		Path:    "/" + dbConfig.Database,
		RawPath: "/" + url.PathEscape(dbConfig.Database),
	}
	query := dsn.Query()
	query.Set("connect_timeout", "10")
	query.Set("timezone", "UTC")
	dsn.RawQuery = query.Encode()

	var err error
	db, err = gorm.Open(
		postgres.Open(dsn.String()),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger: &gormZapLogger{
				logLevel:                  parseLogLevel(dbConfig.LogLevel),
				slowThreshold:             dbConfig.SlowThreshold,
				ignoreRecordNotFoundError: dbConfig.IgnoreRecordNotFoundError,
			},
		},
	)
	if err != nil {
		log.Fatalf("[PostgreSQL] init connection failed: %v\n", err)
	}

	// Trace 注入，不在追踪中展开 SQL 绑定参数。
	if err := db.Use(tracing.NewPlugin(
		tracing.WithoutMetrics(),
		tracing.WithoutQueryVariables(),
		tracing.WithDBSystem("postgresql"),
		tracing.WithAttributes(
			attribute.String("db.namespace", dbConfig.Database),
			attribute.String("server.address", dbConfig.Host),
			attribute.Int("server.port", dbConfig.Port),
		),
	)); err != nil {
		log.Fatalf("[PostgreSQL] init trace failed: %v\n", err)
	}

	// 获取通用数据库对象。
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[PostgreSQL] load sql db failed: %v\n", err)
	}

	// 设置连接池参数。
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Second)
}

// DB 返回绑定请求上下文的 GORM 会话，仅在数据库启用时调用。
func DB(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

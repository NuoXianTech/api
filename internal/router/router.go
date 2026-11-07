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

package router

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/NuoXianTech/api/docs"
	"github.com/NuoXianTech/api/internal/apps/health"
	"github.com/NuoXianTech/api/internal/config"
	"github.com/NuoXianTech/api/internal/db"
	"github.com/NuoXianTech/api/internal/otel_trace"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Serve() {
	var serveErr error
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := otel_trace.Shutdown(ctx); err != nil {
			serveErr = errors.Join(serveErr, fmt.Errorf("shutdown trace: %w", err))
		}
		if serveErr != nil {
			log.Fatalf("[API] serve api failed: %v", serveErr)
		}
	}()
	serveErr = serve()
}

func serve() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if config.Config.Database.Enabled {
		sqlDB, err := db.DB(ctx).DB()
		if err != nil {
			return err
		}
		defer func() {
			if err := sqlDB.Close(); err != nil {
				log.Printf("[PostgreSQL] close connection failed: %v", err)
			}
		}()
	}

	// 运行模式
	if config.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化路由
	r := gin.New()
	r.Use(gin.Recovery())

	// 补充中间件
	r.Use(otelgin.Middleware(config.Config.App.AppName), loggerMiddleware())

	apiGroup := r.Group(config.Config.App.APIPrefix)
	{
		if config.Config.App.Env == "development" {
			// Swagger
			apiGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}

		// API V1
		apiV1Router := apiGroup.Group("/v1")
		{
			// Health
			apiV1Router.GET("/health", health.Health)
			apiV1Router.GET("/ready", health.Ready)
		}
	}

	// Serve
	server := &http.Server{
		Addr:              config.Config.App.Addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return serveHTTP(ctx, server)
}

func serveHTTP(ctx context.Context, server *http.Server) error {
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	// 停止接收请求，等待正在处理的请求完成
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return errors.Join(err, server.Close())
	}
	if err := <-serverErrors; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

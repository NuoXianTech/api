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

package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	ErrorMsg string      `json:"error_msg"`
	Data     interface{} `json:"data"`
}

// Health 检查 HTTP 服务是否存活。
// @Summary 健康检查
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{})
}

type ReadyResponse struct {
	ErrorMsg string      `json:"error_msg"`
	Data     interface{} `json:"data"`
}

// Ready godoc
// @Summary 就绪检查
// @Description 当前仅检查 HTTP 服务；数据库和 Redis 检查将在对应模块接入后补齐。
// @Tags health
// @Produce json
// @Success 200 {object} ReadyResponse
// @Router /ready [get]
func Ready(c *gin.Context) {
	// TODO: 接入 internal/db 后，使用请求上下文检查 PostgreSQL 和 Redis 连接。
	c.JSON(http.StatusOK, ReadyResponse{})
}

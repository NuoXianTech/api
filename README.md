# API

API是一个API 发布、访问控制、调用统计、积分计费与运营的接口管理平台。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.26.10-blue.svg)](https://go.dev/)
[![Next.js](https://img.shields.io/badge/Next.js-16-black.svg)](https://nextjs.org/)
[![React](https://img.shields.io/badge/React-19-blue.svg)](https://react.dev/)

[![GitHub release](https://img.shields.io/github/v/release/NuoXianTech/api?include_prereleases)](https://github.com/NuoXianTech/api/releases)
[![GitHub stars](https://img.shields.io/github/stars/NuoXianTech/api)](https://github.com/NuoXianTech/api/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/NuoXianTech/api)](https://github.com/NuoXianTech/api/network)
[![GitHub issues](https://img.shields.io/github/issues/NuoXianTech/api)](https://github.com/NuoXianTech/api/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/NuoXianTech/api)](https://github.com/NuoXianTech/api/pulls)
[![GitHub contributors](https://img.shields.io/github/contributors/NuoXianTech/api)](https://github.com/NuoXianTech/api/graphs/contributors)

[![Backend Build](https://github.com/NuoXianTech/api/actions/workflows/build_backend.yml/badge.svg)](https://github.com/NuoXianTech/api/actions/workflows/build_backend.yml)
[![Frontend Build](https://github.com/NuoXianTech/api/actions/workflows/build_frontend.yml/badge.svg)](https://github.com/NuoXianTech/api/actions/workflows/build_frontend.yml)
[![CodeQL](https://github.com/NuoXianTech/api/actions/workflows/codeql.yml/badge.svg)](https://github.com/NuoXianTech/api/actions/workflows/codeql.yml)
[![ESLint](https://github.com/NuoXianTech/api/actions/workflows/eslint.yml/badge.svg)](https://github.com/NuoXianTech/api/actions/workflows/eslint.yml)

## 📖 项目简介

**api** 是一个API 发布、访问控制、调用统计、积分计费与运营的接口管理平台。

### ✨ 主要特性

- xxx

## 🏗️ 架构概览

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │    │   Database      │
│   (Next.js)     │◄──►│     (Go)        │◄──►│  (PGSQL/Redis)  │
│                 │    │                 │    │                 │
│ • React 19      │    │ • Gin Framework │    │ • PostgreSQL    │
│ • TypeScript    │    │ • Worker        │    │ • Redis Cache   │
│ • Tailwind CSS  │    │ • OpenTelemetry │    │ • Session Store │
│ • Shadcn UI     │    │ • Swagger API   │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🛠️ 技术栈

### 后端

- **[Go 1.26.10](https://go.dev/doc)** - 主要开发语言
- **[Gin](https://github.com/gin-gonic/gin)** - Web 框架
- **[PostgreSQL](https://www.postgresql.org/)** - 主数据库
- **[Redis](https://github.com/redis/redis)** - 缓存和会话存储
- **[OpenTelemetry](https://opentelemetry.io)** - 可观测性
- **[Swagger](https://github.com/swaggo/swag)** - API 文档

### 前端

- **[Next.js 16](https://github.com/vercel/next.js)** - React 框架
- **[React 19](https://github.com/facebook/react)** - UI 库
- **[TypeScript](https://github.com/microsoft/TypeScript)** - 类型安全
- **[Tailwind CSS 4](https://github.com/tailwindlabs/tailwindcss)** - 样式框架
- **[Shadcn UI](https://github.com/shadcn-ui/ui)** - 组件库
- **[Lucide Icons](https://github.com/lucide-icons/lucide)** - 图标库

## 📋 环境要求

- **Go** 1.26.10
- **Node.js** 24.22.1
- **PostgreSQL** >= 18.6
- **Redis** >= 8.0
- **pnpm** >= 8.0

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/NuoXianTech/api.git
cd api
```

### 2. 配置环境

```bash
cp config.example.yaml config.yaml
```

编辑 `config.yaml` 文件，配置数据库连接、Redis等信息。

### 3. 初始化数据库

```bash
# 创建数据库
psql -U postgres -c "CREATE DATABASE api WITH ENCODING 'UTF8';"

# 运行迁移（启动后端时会自动执行）
```

### 4. 启动后端

```bash
# 安装 Go 依赖
go mod tidy

# 生成 API 文档
make swagger

# 启动后端服务
go run main.go api
```

### 5. 启动前端

```bash
cd frontend

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev
```

### 6. 访问应用

- **前端界面**: http://localhost:3000
- **API 文档**: http://localhost:8000/api/swagger/index.html
- **健康检查**: http://localhost:8000/api/v1/health
- **就绪检查**: http://localhost:8000/api/v1/ready

## ⚙️ 配置说明

### 主要配置项

| 配置项 | 说明 | 示例 |
|--------|------|------|
| `app.addr` | 后端服务监听地址 | `:8000` |
| `worker.port` | Worker 探针端口 | `8001` |
| `database.host` | PostgreSQL 数据库地址 | `127.0.0.1` |
| `redis.host` | Redis 服务器地址 | `127.0.0.1` |

详细配置说明请参考 `config.example.yaml` 文件。

## 🔧 开发指南

### 后端开发

```bash
# 运行 API 服务器
go run main.go api

# 运行工作队列
go run main.go worker

# 生成 Swagger 文档
make swagger

# 代码格式化和检查
make tidy
```

### 前端开发

```bash
cd frontend

# 开发模式（使用 Turbopack）
pnpm dev

# 构建生产版本
pnpm build

# 启动生产服务
pnpm start

# 代码检查和格式化
pnpm lint
pnpm format
```

## 📚 API 文档

API 文档通过 Swagger 自动生成，启动后端服务后可访问：

```
http://localhost:8000/api/swagger/index.html
```

### 主要 API 端点

- `GET /api/v1/health` - 健康检查
- `GET /api/v1/ready` - 就绪检查

## 🧪 测试

```bash
# 后端测试
go test ./...

# 前端测试
cd frontend
pnpm test
```

## 🚀 部署

### Docker 部署

```bash
# 构建镜像
docker build -t api .

# 运行容器
docker run -d -p 8000:8000 api
```

### 生产环境部署

1. 构建前端资源：
   ```bash
   cd frontend && pnpm build
   ```

2. 编译后端程序：
   ```bash
   go build -o api main.go
   ```

3. 配置生产环境的 `config.yaml`

4. 启动服务：
   ```bash
   ./api api
   ```

## 🤝 贡献指南

我们欢迎社区贡献！请在提交代码前阅读：

- [贡献指南](CONTRIBUTING.md)
- [行为准则](CODE_OF_CONDUCT.md)
- [贡献者许可协议](CLA.md)

### 提交流程

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/your-feature`)
3. 提交更改 (`git commit -am 'Add your feature'`)
4. 推送到分支 (`git push origin feature/your-feature`)
5. 创建 Pull Request

提交信息遵循 `<type>(<scope>): <subject>` 格式，且**必须包含 Scope**。

## 📄 许可证

本项目基于 [MIT 许可证](LICENSE) 开源。

## 🔗 相关链接

- [问题反馈](https://github.com/NuoXianTech/api/issues)

## ❤️ 致谢

感谢所有为本项目做出贡献的开发者，以及参与玩法数值与平衡性测试的社区成员。

## 📈 项目趋势

[![Star History Chart](https://api.star-history.com/svg?repos=NuoXianTech/api&type=Date)](https://star-history.com/#NuoXianTech/api&Date)

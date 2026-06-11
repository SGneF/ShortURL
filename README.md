# shortURL

基于 Go 的高性能短链接服务，使用 Gin + GORM + Redis + Bloom Filter 构建。

## 功能

- **短链接生成** — 长链接通过 MD5 去重，自增序列 ID 经 Base62 编码生成短码
- **重定向** — 访问短链接自动跳转至原始 URL
- **布隆过滤器** — 快速排除不存在的短码，避免无效请求穿透到数据库
- **Redis 缓存** — 7 天 TTL 缓存热点链接，减少 DB 压力
- **singleflight** — 并发请求合并，防止缓存击穿
- **敏感词过滤** — 自动跳过生成的敏感短码

## 项目结构

```
shortURL/
├── config/          # 配置加载（YAML + 默认值）
├── dao/             # 数据访问层（MySQL / Redis）
├── handler/         # HTTP 请求处理
├── model/           # 数据模型（GORM）
├── pkg/
│   ├── base62/      # Base62 编解码
│   ├── bloom/       # 布隆过滤器封装
│   └── filter/      # 敏感词过滤
├── router/          # 路由注册
├── service/         # 业务逻辑层
└── main.go          # 入口
```

## API

### 生成短链接

```
POST /api/shorten
Content-Type: application/json

{"long_url": "https://example.com/very/long/url"}
```

响应：

```json
{
  "short_url": "http://localhost:8080/1a2B3c",
  "long_url": "https://example.com/very/long/url"
}
```

### 访问短链接

```
GET /:shortCode
```

重定向到原始长链接（HTTP 302）。

## 技术架构

```
请求 → Router(Gin) → Handler → Service → DAO(MySQL/Redis)
                              ↓
                         Bloom Filter（快速负向判断）
```

**缩短流程**：MD5 长链接去重 → 自增序列 ID → Base62 编码 → 敏感词检查 → 写入 MySQL

**重定向流程**：Bloom Filter 快速否定 → Redis 缓存命中 → MySQL 回源（singleflight 合并）→ 异步回填缓存

## 快速开始

### 前置依赖

- Go 1.26+
- MySQL 8.0+
- Redis 7.0+

### 1. 创建数据库

```sql
CREATE DATABASE shorturl CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

表结构会由 GORM AutoMigrate 自动创建。

### 2. 配置

复制配置模板并修改数据库和 Redis 连接信息：

```yaml
server:
  port: 8080
  domain: "http://localhost:8080"

mysql:
  dsn: "root:password@tcp(127.0.0.1:3306)/shorturl?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0

bloom:
  capacity: 1000000
  false_positive: 0.001
```

也可通过环境变量指定配置文件路径：

```bash
export CONFIG_PATH=/path/to/config.yaml
```

### 3. 运行

```bash
go run main.go
```

服务启动后访问 `http://localhost:8080`。

### 4. 编译

```bash
go build -o shorturl .
./shorturl
```

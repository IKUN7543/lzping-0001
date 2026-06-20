# Go-Zero 微服务电商系统

## 项目简介

基于 Go-Zero 框架从零搭建的微服务电商系统，解决传统单体系统扩展性差、高并发库存超卖、商品查询性能瓶颈等问题。通过微服务拆分 + 多种中间件组合方案，实现高性能、高可用、高扩展的电商后端架构。

## 架构设计

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              Client                                     │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         API Gateway / REST API                          │
│  ┌──────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐             │
│  │ User API │  │ Product API│ │  Stock API│ │  Order API│             │
│  └────┬─────┘  └─────┬─────┘  └─────┬─────┘  └─────┬─────┘             │
└───────┼──────────────┼──────────────┼──────────────┼───────────────────┘
        │              │              │              │
        ▼              ▼              ▼              ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                         gRPC Services (ETCD 注册发现)                    │
│  ┌──────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐             │
│  │ User RPC │  │ Product RPC│ │  Stock RPC│ │  Order RPC│             │
│  └────┬─────┘  └─────┬─────┘  └─────┬─────┘  └────┬┬─────┘             │
└───────┼──────────────┼──────────────┼──────────────┼┼───────────────────┘
        │              │              │              │ │
        ▼              ▼              ▼              │ │
┌──────────────┐  ┌─────────┐  ┌───────────┐         │ │     ┌─────────┐
│    MySQL     │  │  Redis  │  │    ES     │◄────────┘ │     │  Kafka  │
│  (主数据存储) │  │ (缓存+锁)│  │ (全文检索) │           │     │(异步削峰)│
└──────────────┘  └─────────┘  └───────────┘           │     └────┬────┘
                                                         │          │
                                                         └──────────┘
```

## 技术栈

| 分类 | 技术 | 版本 | 说明 |
|------|------|------|------|
| 编程语言 | Go | 1.22 | 高效并发、静态强类型 |
| 微服务框架 | go-zero | 1.7.0 | RPC + REST + 服务治理 |
| 数据库 | MySQL | 8.0 | InnoDB 事务 + 行锁 |
| 缓存 | Redis | 7.x | 缓存层、分布式锁、布隆过滤器 |
| 消息队列 | Kafka | 3.6 | 订单异步削峰填谷 |
| 服务发现 | ETCD | 3.5 | 服务注册/发现、配置 |
| 搜索引擎 | ElasticSearch | 8.11 | 商品全文检索、多条件筛选 |
| 认证授权 | JWT | - | 无状态 Token 鉴权 |
| 容器化 | Docker + Compose | - | 一键部署、环境一致 |
| ORM | GORM | 1.25.x | 类型安全、关联查询 |

## 服务划分

### 1. 用户服务 (User Service)
- **端口**: API `8000`, RPC `8001`
- **功能**: 
  - 用户注册 (Bcrypt 密码加密)
  - 用户登录 (JWT Access + Refresh Token)
  - 用户信息查询 / 修改
- **安全**: Bcrypt 哈希 + JWT 双 Token

### 2. 商品服务 (Product Service)
- **端口**: API `8010`, RPC `8011`
- **功能**:
  - 商品 CRUD
  - 商品详情 (**Redis Cache Aside 模式**)
  - 商品列表 / 全文检索 (**ElasticSearch**)
  - 分类树查询
- **性能优化**:
  - **布隆过滤器** → 防缓存穿透
  - **Cache Aside** → 双写一致性
  - **ES 全文检索** → 响应 < 100ms

### 3. 库存服务 (Stock Service)
- **端口**: API `8020`, RPC `8021`
- **功能**:
  - 库存创建 / 查询
  - 库存扣减 / 返还 / 确认
- **防超卖方案**:
  - **Redis 分布式锁** (Lua 原子脚本)
  - **数据库乐观锁** (Version 版本号)
  - **Redis 预扣库存** + DB 异步同步

### 4. 订单服务 (Order Service)
- **端口**: API `8030`, RPC `8031`
- **功能**:
  - 创建订单 (调用库存扣减)
  - 订单查询 / 列表
  - 取消订单 (库存回滚)
  - 订单支付 (库存确认)
- **高并发优化**:
  - **雪花算法** 生成全局唯一订单 ID
  - **Kafka** 异步写入订单 → 削峰填谷
  - **RPC 调用链** → 库存 → 商品

## 项目结构

```
Go-Zero/
├── common/                     # 公共库
│   ├── errx/                   # 错误码定义
│   ├── snowflake/              # 雪花算法 ID
│   ├── bloom/                  # 布隆过滤器
│   ├── redislock/              # Redis 分布式锁
│   ├── cache/                  # Cache Aside 封装
│   ├── kafkaq/                 # Kafka Producer/Consumer
│   └── middleware/             # 中间件 (Prometheus 等)
├── service/
│   ├── user/
│   │   ├── api/                # REST API
│   │   └── rpc/                # gRPC 服务
│   ├── product/
│   │   ├── api/
│   │   └── rpc/
│   ├── stock/
│   │   ├── api/
│   │   └── rpc/
│   └── order/
│       ├── api/
│       └── rpc/
├── deploy/
│   ├── dockerfile/             # 各服务 Dockerfile
│   └── sql/                    # 初始化 SQL
├── docker-compose.yml          # 一键编排
├── go.mod
└── README.md
```

> **设计规范**: 每个服务内部分层为 `api` (HTTP 网关) + `rpc` (gRPC 业务)，服务间通过 ETCD 做服务发现。API 层做 JWT 鉴权、参数校验、RPC 调用；RPC 层做核心业务逻辑、数据访问。

## 快速开始

### 环境要求

- Docker 20+
- Docker Compose v2+
- 至少 6GB 内存 (ES + Kafka 占用较多)

### 一键启动

```bash
# 克隆项目后进入根目录
cd Go-Zero

# 启动全部服务 (首次需构建镜像，耗时较久)
docker-compose up -d

# 查看所有服务状态
docker-compose ps

# 查看某服务日志
docker-compose logs -f order-rpc

# 停止全部
docker-compose down
```

### 启动顺序

Compose 通过 `depends_on` + `healthcheck` 自动控制启动顺序：

```
MySQL, Redis, ETCD, ZooKeeper → Kafka → ElasticSearch 
  → 各 RPC 服务 → 各 API 服务
```

### 端口速查

| 服务 | 端口 | 说明 |
|------|------|------|
| user-api | 8000 | 用户 REST API |
| product-api | 8010 | 商品 REST API |
| stock-api | 8020 | 库存 REST API |
| order-api | 8030 | 订单 REST API |
| MySQL | 3306 | root / 123456 |
| Redis | 6379 | 无密码 |
| ETCD | 2379 | 无密码 |
| Kafka | 9092 | |
| ElasticSearch | 9200 | |
| Kibana | 5601 | ES 可视化控制台 |

## API 使用示例

### 1. 用户注册

```bash
curl -X POST http://localhost:8000/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "zhangsan",
    "password": "123456",
    "nickname": "张三",
    "mobile": "13900001234"
  }'
```

### 2. 用户登录 (获取 Token)

```bash
curl -X POST http://localhost:8000/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'

# 返回:
# {
#   "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
#   "expiresAt": 1718xxx,
#   "id": 2,
#   ...
# }
```

### 3. 查询商品详情

```bash
# 未缓存首次查 DB，平均 200ms；之后查 Redis，<30ms
curl "http://localhost:8010/product/detail?id=1"
```

### 4. 商品全文检索

```bash
# 通过 ES 全文搜索，响应 <100ms
curl "http://localhost:8010/product/search?keyword=iPhone&page=1&pageSize=10"
```

### 5. 初始化库存 (测试账号)

```bash
# 先登录获取 token
TOKEN="粘贴上面的 accessToken"

curl -X POST http://localhost:8020/stock/admin/create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"productId":1,"num":100}'
```

### 6. 创建订单

```bash
curl -X POST http://localhost:8030/order/create \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [{"productId":1,"num":1,"price":999900}],
    "receiverName": "张三",
    "receiverPhone": "13900001234",
    "receiverAddress": "北京市朝阳区xxx"
  }'
```

## 核心技术方案说明

### 🔴 库存超卖防护 (三重保障)

```
请求 → Redis 分布式锁 (原子性，10s 超时)
       → Redis DECRBY 预扣库存 (<0 回滚 + 返回)
       → DB 乐观锁 (WHERE version=xx AND available>=n)
       → 成功: 写幂等键防重复
```

Lua 脚本保证 SETNX + EXPIRE 原子；Redis + DB 双重扣减保证极端一致性。

### 🟠 商品缓存体系

```
查询: 
  BloomFilter 不存在 → 直接返回 (防穿透)
  → Redis Get 命中 → 返回
  → DB 查询 → 写 Redis (5min TTL) → 返回
更新/删除:
  → DB 写 → Redis DEL → ES 同步
```

*热点数据*: 高销量商品通过布隆过滤器预热进 Redis，命中率 > 90%。

### 🟡 订单异步削峰

```
API 创建订单 → 同步扣库存 → Kafka Send (Topic: order-create-topic)
                                              ↓
Consumer (多实例) → 批量写 DB → 订单持久化
```

相比同步写 DB，峰值吞吐量提升 5~10 倍，削平高并发毛刺。

### 🟢 ES 全文检索

```json
{
  "query": {
    "bool": {
      "must": [
        {"multi_match": {"query": "iPhone", "fields": ["name^3","brand^2","subtitle","detail"]}}
      ],
      "filter": [{"term": {"categoryId": 9}}, {"term": {"status": 1}}]
    }
  }
}
```

字段权重：名称(3) > 品牌(2) > 副标题 > 详情。商品创建/更新自动同步索引。

### 🔵 雪花 ID

1bit 符号 + 41bit 毫秒时间戳 + 10bit 工作节点 + 12bit 序列号，单节点每秒 409万 ID。

## 性能指标 (设计目标)

| 指标 | 目标值 | 达成方案 |
|------|--------|----------|
| 商品详情平均响应 | < 30ms | Redis 缓存 + Cache Aside |
| 商品搜索响应 | < 100ms | ES 倒排索引 + 字段权重 |
| 缓存命中率 | > 90% | 布隆过滤器预热 + 热点数据 |
| 订单峰值 TPS | 提升 5x+ | Kafka 异步削峰 + 幂等 |
| 库存一致性 | 0 超卖 | 分布式锁 + 乐观锁 + 幂等键 |

## 本地开发模式

如果不使用 Docker，可按以下步骤本地启动：

```bash
# 1. 启动中间件 (使用 docker 只起中间件)
docker-compose up -d mysql redis etcd zookeeper kafka elasticsearch

# 2. 修改各服务 etc/*.yaml，把主机名改成 localhost
#    例如: mysql:3306 → 127.0.0.1:3306

# 3. 按顺序启动 RPC
go run service/user/rpc/user.go
go run service/product/rpc/product.go
go run service/stock/rpc/stock.go
go run service/order/rpc/order.go

# 4. 启动 API
go run service/user/api/user.go
go run service/product/api/product.go
go run service/stock/api/stock.go
go run service/order/api/order.go
```

## 生产环境建议

1. **MySQL**: 主从复制 + 分库分表 (订单按月分表)
2. **Redis**: 哨兵/集群模式 + 读写分离
3. **Kafka**: 多分区 + 多 Consumer Group + 死信队列
4. **ETCD**: 3 节点集群，启用 TLS 认证
5. **ES**: 冷热分离 + 分片副本 + Snapshot 备份
6. **全链路**: OpenTelemetry 埋点 + Jaeger 追踪
7. **监控**: Prometheus + Grafana + 告警规则
8. **网关**: 接入 Nginx/APISIX 做统一入口、限流、熔断

## License

MIT

---

**如果这个项目对你有帮助，欢迎 Star 🌟**

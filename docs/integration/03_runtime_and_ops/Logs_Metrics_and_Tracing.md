# 插件日志、指标与链路追踪规范（03_runtime_and_ops/Logs_Metrics_and_Tracing.md）

> 本文档定义 PowerX 插件的运行日志（Logs）、性能指标（Metrics）与链路追踪（Tracing）规范，  
> 以统一插件与宿主的可观测性标准，实现跨模块、跨租户的全栈监控与分析。

---

## 🧭 一、文档目标

- 统一 PowerX 插件的日志与监控体系；
- 定义日志格式、等级与落盘规则；
- 说明插件侧 Metrics（指标）上报方式；
- 建立宿主与插件的链路追踪（Trace Context）传播标准；
- 支持日志采集、Grafana 监控、Jaeger/Tempo 链路分析。

---

## 🧱 二、整体监控架构

```

┌────────────────────────────────────────────┐
│              PowerX Observability           │
│ ┌────────────┐  ┌────────────┐  ┌─────────┐ │
│ │   Loki     │  │ Prometheus │  │  Tempo  │ │
│ └────────────┘  └────────────┘  └─────────┘ │
│      ↑                 ↑               ↑     │
│ ┌────────────┐  ┌────────────┐  ┌────────────┐│
│ │ Plugin CRM │  │ Plugin EC  │  │ Plugin AI  ││
│ └────────────┘  └────────────┘  └────────────┘│
│              ↑                 │               │
│         PowerX Core (Agent + Runtime Manager)  │
└────────────────────────────────────────────┘

```

---

## 🧩 三、日志系统（Logs）

### 1️⃣ 日志文件目录

插件运行时日志默认存放在：

```

/var/lib/powerx/plugins/<plugin_id>/logs/
├── runtime.log     # 主运行日志
├── access.log      # HTTP/gRPC 访问日志
├── error.log       # 错误日志
├── metrics.log     # 指标输出（文本格式）
└── trace.log       # Trace 链路采样

```

### 2️⃣ 日志级别与格式

| 等级 | 用途 | 示例 |
|------|------|------|
| `DEBUG` | 调试信息 | 连接请求、参数值 |
| `INFO` | 常规事件 | 插件启动成功 |
| `WARN` | 非致命警告 | API 限流 |
| `ERROR` | 错误事件 | DB 查询失败 |
| `FATAL` | 致命异常 | 无法启动 |

**推荐统一 JSON 格式输出：**

```json
{
  "timestamp": "2025-10-13T10:32:21Z",
  "level": "INFO",
  "plugin_id": "com.powerx.plugin.crm",
  "tenant_id": "tenant_123",
  "component": "contact_service",
  "message": "Contact created successfully",
  "trace_id": "2cf0c83d8bfc49cb"
}
```

> PowerX Core 的日志收集器会自动识别 `plugin_id` 与 `trace_id` 字段。
> 插件可通过 stdout 输出日志，由宿主捕获并写入文件。

---

## 🧩 四、日志采集与上报

宿主的 `powerx-log-agent` 进程会周期性（默认 60s）扫描：

```
/var/lib/powerx/plugins/*/logs/*.log
```

并上传到 Loki / Elasticsearch。

配置示例：

```yaml
log_agent:
  upload_interval: 60s
  include:
    - runtime.log
    - error.log
  exclude:
    - metrics.log
  target: loki://loki.powerx.cloud
```

插件开发者可通过标签区分日志：

```
[INFO] [plugin=crm] [tenant=123] [trace=abc123] created contact id=456
```

---

## ⚙️ 五、Metrics（指标系统）

插件应暴露 Prometheus 格式指标端点：

```
GET /metrics
→ Content-Type: text/plain; version=0.0.4
```

### 示例输出

```
# HELP plugin_request_total Total requests handled by plugin
# TYPE plugin_request_total counter
plugin_request_total{plugin="crm", capability="contact.create"} 1532

# HELP plugin_latency_seconds Request latency
# TYPE plugin_latency_seconds histogram
plugin_latency_seconds_bucket{le="0.1"} 300
plugin_latency_seconds_bucket{le="0.5"} 820
plugin_latency_seconds_sum 90.3
plugin_latency_seconds_count 1020
```

### 必须指标

| 指标名                        | 类型        | 说明         |
| -------------------------- | --------- | ---------- |
| `plugin_request_total`     | counter   | 插件处理的总请求数  |
| `plugin_error_total`       | counter   | 错误请求总数     |
| `plugin_latency_seconds`   | histogram | 请求耗时分布     |
| `plugin_cpu_seconds_total` | counter   | CPU 使用量    |
| `plugin_memory_bytes`      | gauge     | 内存占用       |
| `plugin_mcp_sessions`      | gauge     | 当前 MCP 连接数 |

宿主统一抓取：

```
http://127.0.0.1:<PORT>/metrics
```

> 插件必须在 manifest 中声明 metrics 端点（见 [Runtime_Env_and_Ports.md](./Runtime_Env_and_Ports.md)）。

---

## 🧠 六、Tracing（链路追踪）

PowerX 采用 OpenTelemetry 协议实现插件链路追踪。

### 1️⃣ Trace Context 传播

宿主在调用插件能力时会注入请求头：

```
traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
```

插件需在处理逻辑中：

- 提取 Trace Context；
- 生成子 Span；
- 将上下文传递给下游调用（如 DB / 外部 API）。

Go Gin 示例：

```go
func MiddlewareTracing() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
        spanCtx, span := tracer.Start(ctx, c.FullPath())
        defer span.End()
        c.Request = c.Request.WithContext(spanCtx)
        c.Next()
    }
}
```

### 2️⃣ 插件输出 Span

每个能力调用应生成 Span：

```
powerx.plugin.crm.contact.create
```

属性：

| Key                | 示例                    | 说明   |
| ------------------ | --------------------- | ---- |
| `plugin.id`        | com.powerx.plugin.crm | 插件标识 |
| `tenant.id`        | tenant_123            | 租户   |
| `capability`       | crm.contact.create    | 能力名称 |
| `http.status_code` | 200                   | 状态码  |
| `latency.ms`       | 42                    | 延迟   |

---

## 🧩 七、Tracing 收集与可视化

宿主 Core 内置 OpenTelemetry Collector，将数据导出至 Tempo / Jaeger。

配置示例：

```yaml
tracing:
  exporter: otlp
  endpoint: tempo.powerx.cloud:4317
  sampling_rate: 0.1
```

Grafana 展示路径：

```
PowerX → Observability → Plugin Traces
```

---

## 🧩 八、Metrics 与 Logs 的关联追踪

每个日志条目必须包含 Trace ID：

```
{"trace_id": "4bf92f3577b34da6a3ce929d0e0e4736"}
```

通过 Trace ID，Grafana 可直接关联：

- 请求详情（Logs）
- 性能统计（Metrics）
- 调用链图（Traces）

---

## 🧱 九、插件侧日志封装建议（Go 实现）

推荐使用 `zap` 或 `logrus`：

```go
import "go.uber.org/zap"

var logger, _ = zap.NewProduction()

func main() {
    logger.Info("plugin started",
        zap.String("plugin_id", os.Getenv("POWERX_PLUGIN_ID")),
        zap.String("version", os.Getenv("POWERX_PLUGIN_VERSION")),
        zap.Time("time", time.Now()),
    )
}
```

日志示例：

```
{"level":"info","plugin_id":"com.powerx.plugin.crm","msg":"plugin started","time":"2025-10-13T09:00:00Z"}
```

---

## 🧩 十、日志与指标的采样策略

| 类型      | 采样率             | 说明             |
| ------- | --------------- | -------------- |
| Trace   | 10% 默认采样        | 高延迟请求强制采样      |
| Metrics | 全量统计            | 1s 刷新周期        |
| Logs    | 全量记录 WARN+ERROR | DEBUG 仅在开发环境开启 |

---

## 🧩 十一、运维与报警集成

宿主通过 Prometheus Alertmanager 管理报警：

| 指标                       | 触发条件         | 动作       |
| ------------------------ | ------------ | -------- |
| `plugin_error_total`     | 5min 内 > 100 | 邮件通知     |
| `plugin_latency_seconds` | P95 > 1s     | Slack 告警 |
| `powerx_mcp_sessions`    | 断开数 > 10     | 重启插件     |
| `quota_exceeded_total`   | > 0          | 禁用调用     |

---

## ⚙️ 十二、开发与测试建议

- 使用 `POWERX_ENV=dev` 输出详细日志；
- 在本地测试 `/metrics` 与 `/healthz`；
- 在 CI 中集成日志与指标验证：

  ```bash
  make test-metrics
  ```

- 检查 Trace Context 是否正确传递至下游；
- 日志输出不得包含敏感数据（Token、密码、邮箱）。

---

## 🧩 十三、安全与合规要求

| 项目         | 要求                    |
| ---------- | --------------------- |
| **隐私数据脱敏** | 日志中禁止直接输出手机号、邮箱、Token |
| **多租户隔离**  | 必须在日志中显式标注 tenant_id  |
| **日志保留策略** | 默认 30 天，敏感插件可自定义      |
| **安全访问**   | 所有指标端点需鉴权或仅限宿主访问      |

---

## 📚 延伸阅读

- [Runtime_Env_and_Ports.md](./Runtime_Env_and_Ports.md)
- [MCP_Session_and_Registration.md](./MCP_Session_and_Registration.md)
- [Quotas_RateLimits_and_Costs.md](./Quotas_RateLimits_and_Costs.md)
- [04_security_and_compliance/Plugin_Security_Checklist.md](../04_security_and_compliance/Plugin_Security_Checklist.md)

---

> **文档版本：** v1.0.0
> **适用范围：** PowerX ≥ 0.9.0
> **维护团队：** PluginBase Observability Team
> **最后更新：** 2025-10

```

---

✅ **总结要点：**

- 定义了插件日志输出、指标暴露与链路追踪的标准；
- 全面兼容 PowerX 宿主的 Loki / Prometheus / Tempo 三栈；
- 要求日志含 `plugin_id`、`tenant_id`、`trace_id`；
- 支持跨租户可观测性与快速调试；
- 下一章节将进入 **04_security_and_compliance**，开始定义插件安全、隐私与漏洞响应规范。  

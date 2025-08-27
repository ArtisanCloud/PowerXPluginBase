# 数据库连接说明

## 数据库信息

- **数据库名**: `powerx_plugin_scrum`
- **用户**: `michaelhu` (拥有完整权限)
- **Schema**: `scrum`
- **端口**: 5432 (PostgreSQL 默认)

## 连接方式

### 1. 命令行连接
```bash
psql -d powerx_plugin_scrum
```

### 2. 使用完整连接字符串
```bash
psql "host=localhost user=michaelhu dbname=powerx_plugin_scrum port=5432 sslmode=disable"
```

### 3. 应用程序连接 (DSN)
```
host=localhost user=michaelhu dbname=powerx_plugin_scrum port=5432 sslmode=disable TimeZone=Asia/Shanghai
```

## 数据库结构

### 表结构
- **scrum.sprint**: Sprint 管理表
- **scrum.task**: 任务管理表

### 示例数据
数据库已包含示例数据：
- 2 个租户 (tenant_id: 1, 2)
- 每租户 3 个 Sprint
- 每租户 12 个任务（包含不同状态和优先级）

### 查看数据
```sql
-- 切换到 scrum schema
SET search_path TO scrum;

-- 查看所有 Sprint
SELECT * FROM sprint ORDER BY tenant_id, id;

-- 查看所有任务
SELECT id, title, status, priority, tenant_id, sprint_id FROM task ORDER BY tenant_id, id;

-- 查看租户1的进行中任务
SELECT * FROM task WHERE tenant_id = 1 AND status = 'in_progress';
```

## 权限说明

用户 `michaelhu` 拥有以下权限：
- ✅ 数据库连接权限
- ✅ Schema 使用和创建权限  
- ✅ 所有表的 CRUD 权限
- ✅ 序列和函数的完整权限
- ✅ 未来创建对象的默认权限

## 配置文件

项目配置文件已更新：
- `backend/etc/config.yaml` - 实际配置（已设置正确的数据库连接）
- `backend/etc/config.example.yaml` - 示例配置

## 测试连接

可以通过以下方式测试数据库连接是否正常：

```bash
# 1. 运行数据库迁移
cd backend
PX_DEV_MODE=true go run ./cmd/database/migrate

# 2. 运行种子数据
PX_DEV_MODE=true go run ./cmd/database/seed

# 3. 启动应用程序
PX_DEV_MODE=true go run ./cmd/plugin
```

数据库现在已经准备就绪，可以开始开发和测试了！🎉
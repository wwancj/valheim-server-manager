# Valheim Server Manager - 重构为纯 Web 服务

## 背景

Wails 绑定存在稳定性问题（返回值序列化失败），决定改为 **Go HTTP Server + React SPA** 架构。

## 架构变更

**之前**: Wails 桌面应用 (Go ↔ WebView2 ↔ React)
**之后**: Go HTTP Server 提供 REST API + WebSocket，React 前端通过浏览器访问

```
Go HTTP Server (:13256)
├── GET  /api/app/info
├── GET  /api/servers
├── POST /api/servers
├── PUT  /api/servers/:id
├── DELETE /api/servers/:id
├── POST /api/servers/:id/start
├── POST /api/servers/stop
├── GET  /api/servers/status
├── GET  /api/steamcmd/status
├── POST /api/steamcmd/install
├── POST /api/valheim/install
├── GET  /api/config/:id
├── PUT  /api/config/:id
├── GET  /api/logs
├── DELETE /api/logs
├── GET  /api/mods/search?query=&page=
├── GET  /api/mods/detail/:namespace/:name
├── GET  /api/mods/installed?serverDir=
├── POST /api/mods/install
├── DELETE /api/mods/installed
├── PUT  /api/mods/toggle
├── GET  /api/system
├── GET  /api/system/process/:pid
├── POST /api/backups
├── GET  /api/backups/:serverId
├── POST /api/backups/:id/restore
├── DELETE /api/backups/:id
├── GET  /api/profiles?serverDir=
├── POST /api/profiles
├── DELETE /api/profiles/:id
├── POST /api/safety/check
├── POST /api/bepinex/install
├── GET  /api/bepinex/status
├── GET  /api/bepinex/plugins
├── WS   /ws/events (实时日志/进度/状态)
└──      / (React SPA 静态文件)
```

## 技术选型

- **HTTP 路由**: `net/http` + `gorilla/mux` (轻量、成熟)
- **WebSocket**: `gorilla/websocket` (实时事件推送)
- **静态文件**: Go `embed` 嵌入 React 构建产物
- **前端请求**: `fetch` API (替换 Wails 绑定调用)
- **实时通信**: WebSocket 替换 Wails Events

## 实施步骤

### Step 1: Go 后端改造
1. 创建 `app/server/router.go` - 路由注册
2. 创建 `app/server/handlers.go` - HTTP 处理函数
3. 创建 `app/server/websocket.go` - WebSocket 事件推送
4. 创建 `app/server/middleware.go` - CORS、日志中间件
5. 修改 `main.go` - 移除 Wails，启动 HTTP Server
6. 更新 `go.mod` - 移除 wails，添加 mux/websocket

### Step 2: 前端改造
1. 创建 `frontend/src/api/client.ts` - HTTP + WebSocket 客户端
2. 创建 `frontend/src/api/types.ts` - 复用现有 models.ts 类型
3. 修改所有页面组件 - 用 fetch 替换 wails 绑定调用
4. 修改实时功能 - 用 WebSocket 替换 wails Events
5. 删除 `frontend/wailsjs/` 目录
6. 更新 `vite.config.ts` - 添加 API 代理配置

### Step 3: 构建与测试
1. `go mod tidy` 更新依赖
2. `npm run build` 构建前端
3. `go run .` 或 `wails build` 测试
4. 验证所有 API 端点
5. 验证实时功能（日志、进度、状态）

## 保留不变

- 所有业务逻辑 (steamcmd, process, config, mods, backup, etc.)
- 数据模型 (models.go)
- 文件存储结构
- UI 设计和布局

## 删除

- `frontend/wailsjs/` (自动生成的绑定)
- `wails.json` (Wails 配置)
- Wails 依赖 (go.mod 中)

## 端口

默认: `13256` (可在 main.go 中修改)

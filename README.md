<div align="center">

# 🎮 Valheim Server Manager

### 英灵神殿服务器管理工具

[English](#english) | 中文

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![React](https://img.shields.io/badge/React-18-61DAFB?logo=react&logoColor=black)](https://reactjs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org)
[![Ant Design](https://img.shields.io/badge/Ant%20Design-5-0170FE?logo=antdesign&logoColor=white)](https://ant.design)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

一款一站式 Valheim 专用服务器管理工具，支持服务器启停、SteamCMD 安装、Mod 管理、配置编辑、备份恢复、实时日志等功能。

</div>

---

## ✨ 功能特性

### 🖥️ 服务器管理
- 服务器创建、启动、停止、删除
- 支持 **Crossplay 开关**（跨平台联机 / IP 直连切换）
- 一键导出 `start_server.bat` 启动脚本
- 多服务器实例管理

### ⚙️ 配置管理
- 可视化配置编辑器（服务器参数、世界规则、难度预设）
- 7 种内置游戏模式预设（原版、休闲、简单、困难、极限、沉浸、建筑）
- 自定义预设保存与加载
- 管理员/白名单/黑名单管理

### 📦 Mod 管理
- **Thunderstore** 在线 Mod 搜索与浏览
- **BepInEx** 一键安装
- Mod 安装、卸载、启用/禁用
- 依赖检测与冲突检查
- Mod 配置文件管理

### 💾 备份与恢复
- 一键创建服务器备份
- 备份列表管理与快速恢复
- 自动备份策略配置

### 📊 监控与日志
- 实时服务器日志（WebSocket 推送）
- 系统资源监控（CPU、内存、磁盘）
- 服务器进程状态监控
- 日志暂停、清空、导出

### 🛡️ 安全检查
- 启动前安全检查（目录、端口、密码、世界文件）
- 端口冲突检测
- 密码强度提示

---

## 🚀 快速开始

### 环境要求

| 依赖 | 版本 |
|------|------|
| Go | 1.21+ |
| Node.js | 18+ |
| npm | 9+ |

### 安装与运行

```bash
# 1. 克隆仓库
git clone https://github.com/wwancj/valheim-server-manager.git
cd valheim-server-manager

# 2. 安装前端依赖
cd frontend
npm install
cd ..

# 3. 安装 Go 依赖
go mod tidy

# 4. 开发模式运行
go run main.go
# 或者构建后运行
go build -o bin/valheim-server-manager.exe .
./bin/valheim-server-manager.exe
```

启动后访问 **http://localhost:13256**

### 使用流程

```
1️⃣  添加服务器（填写名称、安装路径、端口）
2️⃣  安装 SteamCMD → 安装 Valheim Dedicated Server
3️⃣  配置服务器参数（密码、游戏模式、Crossplay）
4️⃣  启动服务器
5️⃣  在 Valheim 客户端连接：
     • Crossplay 关闭 → 输入 IP:端口（如 127.0.0.1:2456）
     • Crossplay 开启 → 通过 Steam 好友邀请加入
```

---

## 📁 项目结构

```
valheim-server-manager/
├── main.go                     # 程序入口，HTTP 服务器启动
├── app/
│   ├── app.go                  # 应用主结构，业务方法绑定
│   ├── server/                 # HTTP 路由、处理器、WebSocket
│   │   ├── router.go           # API 路由定义
│   │   ├── handlers.go         # 请求处理器
│   │   ├── middleware.go        # CORS、日志中间件
│   │   └── websocket.go        # WebSocket 实时推送
│   ├── config/                 # 服务器配置管理
│   ├── process/                # 服务器进程管理
│   ├── steamcmd/               # SteamCMD 安装与管理
│   ├── mods/                   # Mod 安装与管理
│   ├── thunderstore/           # Thunderstore API 客户端
│   ├── backup/                 # 备份与恢复
│   ├── services/               # 业务服务层
│   ├── system/                 # 系统监控
│   ├── events/                 # 事件发射器
│   └── models/                 # 数据模型
├── frontend/
│   ├── src/
│   │   ├── api/                # HTTP API 客户端
│   │   ├── pages/              # 页面组件（7个功能页）
│   │   ├── components/         # 通用组件
│   │   └── contexts/           # React Context
│   └── dist/                   # 构建产物（嵌入 Go 二进制）
└── bin/                        # 编译输出
```

---

## 🛠️ 技术架构

```
┌─────────────────────────────────────┐
│          浏览器 (React SPA)          │
│   React + TypeScript + Ant Design   │
└──────────────┬──────────────────────┘
               │ HTTP / WebSocket
┌──────────────▼──────────────────────┐
│        Go HTTP Server (gorilla)      │
│   REST API + WebSocket 实时推送      │
├─────────────────────────────────────┤
│          业务服务层 (Services)        │
│  Server / SteamCMD / Mod / Backup   │
│  Config / Log / System / Safety     │
├─────────────────────────────────────┤
│          系统层 (OS Operations)      │
│   进程管理 / 文件系统 / 网络请求      │
└─────────────────────────────────────┘
```

### 技术栈

| 层级 | 技术 |
|------|------|
| 桌面/服务端 | Go + gorilla/mux + gorilla/websocket |
| 前端 | React 18 + TypeScript + Ant Design 5 |
| 构建 | Vite (前端) + Go embed (嵌入静态文件) |
| 通信 | REST API + WebSocket (实时事件) |

---

## 🔌 API 端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/servers` | 获取服务器列表 |
| POST | `/api/servers` | 创建服务器 |
| POST | `/api/servers/start-with-config` | 启动服务器（带配置） |
| POST | `/api/servers/stop` | 停止服务器 |
| POST | `/api/servers/export-script` | 导出启动脚本 |
| GET | `/api/config/templates` | 获取配置模板 |
| GET | `/api/mods/search` | 搜索 Mod |
| POST | `/api/mods/install` | 安装 Mod |
| GET | `/api/system` | 系统信息 |
| GET | `/api/ws/events` | WebSocket 事件流 |

---

## 📝 连接说明

### 本机连接（127.0.0.1）

1. 启动服务器时 **关闭 Crossplay**
2. Valheim 客户端 → 加入游戏 → IP 加入
3. 输入 `127.0.0.1:2456`（端口为你配置的端口）

### 远程连接

1. 关闭 Crossplay
2. 路由器端口转发：UDP 2456-2457
3. Windows 防火墙放行 UDP 2456-2457
4. 使用公网 IP:端口 连接

### 跨平台联机（Crossplay）

1. 启动服务器时 **开启 Crossplay**
2. 通过 Steam 好友邀请或社区服务器列表加入
3. 支持 Steam / Xbox 跨平台

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/xxx`)
3. 提交更改 (`git commit -m 'feat: add xxx'`)
4. 推送分支 (`git push origin feature/xxx`)
5. 创建 Pull Request

---

## 📄 许可证

[MIT License](LICENSE)

---

<div align="center">

## English

A desktop application for managing Valheim Dedicated Server, BepInEx mods, and Thunderstore packages.

**Features:** Server lifecycle management · SteamCMD integration · Mod management (BepInEx + Thunderstore) · Visual config editor · Backup & restore · Real-time logs (WebSocket) · System monitoring · Safety checks · Crossplay toggle · Startup script export

**Quick Start:**
```bash
git clone https://github.com/wwancj/valheim-server-manager.git
cd valheim-server-manager/frontend && npm install && cd ..
go mod tidy && go build -o bin/valheim-server-manager.exe .
./bin/valheim-server-manager.exe
# Open http://localhost:13256
```

**Connect locally:** Disable Crossplay → Join Game → Enter `127.0.0.1:2456`

</div>
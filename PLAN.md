# Valheim Server Manager — Coding Plan

## Phase 0：项目骨架 ✅

目标：先跑起来。

实现：
- Wails 初始化
- React + TS
- Ant Design
- 路由
- 左侧导航
- 基础 Layout
- Go → React 调用
- React → Go 调用
- 开发环境启动
- Windows 打包测试

页面：Dashboard / Servers / Mods / Config / Logs / Backup / Settings

验收：启动应用 → 正常打开 → 页面切换正常 → React 可以调用 Go → Go 可以返回数据

---

## Phase 1：服务器目录管理 ⬜

先不要碰 SteamCMD。

实现 Server Model，支持创建/删除/修改/检测服务器文件。

页面：服务器列表

---

## Phase 2：SteamCMD ⬜

第一个真正的核心模块。

实现 SteamCMDManager：检测/下载/安装/执行/读取 stdout/stderr/解析进度。

Valheim Dedicated Server App ID = 896660

---

## Phase 3：服务器启动/停止 ⬜

实现 ServerProcessManager：Start/Stop/Restart/IsRunning/GetPID

启动参数由 ServerConfig + CommandBuilder 生成，不要在 UI 中拼启动命令。

---

## Phase 4：实时日志 ⬜

服务器 stdout/stderr 实时传递：Valheim Process → Go → WebSocket/Event → React

支持：清空、暂停滚动、自动滚动、搜索、错误过滤、导出

---

## Phase 5：服务器配置可视化 ⬜

建立 ServerConfig 接口，UI 可视化配置基础参数和游戏难度。

保存路径：UI → Go → SQLite；启动路径：SQLite → ServerConfig → CommandBuilder → Valheim

---

## Phase 6：BepInEx 管理 ⬜

一键安装 BepInEx。检测 → 下载 → 解压 → 安装 → 检测 plugins/config

Windows/Linux 路径必须抽象，不允许硬编码。

---

## Phase 7：Thunderstore Mod 浏览 ⬜

核心卖点。实现 ThunderstoreClient：搜索/获取信息/版本/下载地址/依赖/manifest/README

UI：Mod 商店搜索、分类、排序、详情、安装

---

## Phase 8：Mod 安装器 ⬜

实现 ModInstaller：manifest.json → 解析依赖 → 建立依赖树 → 下载 → 解压 → 安装 → 验证

必须防止：循环依赖、重复依赖、版本冲突、下载失败、文件覆盖

---

## Phase 9：已安装 Mod ⬜

页面：已安装 Mod 列表，支持启用/禁用/卸载

启用/禁用采用安全的文件策略（.disabled 后缀），不要直接删除。

---

## Phase 10：Mod 配置可视化 ⬜

第一版：文本配置编辑器
第二阶段：CFG 类型解析自动生成表单

---

## Phase 11：服务器监控 ⬜

Dashboard 展示：CPU / Memory / Uptime / Players

---

## Phase 12：自动备份 ⬜

支持：手动备份 / 自动备份 / 恢复 / 删除备份

备份内容：World / Server Config / Mod List / Mod Config

---

## Phase 13：Mod Profile ⬜

支持：创建/复制/导出/导入/切换 Profile

---

## Phase 14：安全机制 ⬜

启动前检查：密码/端口/世界文件/BepInEx/Mod/依赖/配置

发现问题给出可操作建议，而不是直接启动报错。

---

## MVP 最终范围

- [ ] Wails Desktop
- [ ] Valheim Server 安装
- [ ] SteamCMD
- [ ] 一键启动/停止/重启
- [ ] 服务器配置
- [ ] 实时日志
- [ ] BepInEx 安装
- [ ] Thunderstore 搜索
- [ ] Mod 安装/卸载/启用/禁用
- [ ] Mod 依赖
- [ ] Mod 配置
- [ ] 基础服务器监控
- [ ] 世界备份

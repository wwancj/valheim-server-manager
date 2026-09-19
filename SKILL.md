---
name: valheim-dedicated-server
description: Valheim 专用服务器配置与管理技能。用于创建、修改、解释和排查 Valheim Dedicated Server 的启动脚本、服务器参数、世界存档、权限文件、备份、Crossplay、世界修饰器及 Docker 部署。涉及 Valheim 服务器配置时，优先依据官方 Dedicated Server Guide，不虚构不存在的配置文件或参数。
---

# Valheim Dedicated Server Skill

## 目标

处理 Valheim Dedicated Server 的配置、启动、停止、权限管理、存档管理、备份、Crossplay、世界修饰器和 Docker 部署。

核心原则：

1. Valheim Dedicated Server 的核心配置通过启动脚本参数完成。
2. Windows 使用 `start_headless_server.bat`，Linux 使用 `start_server.sh`。
3. 修改官方原始启动脚本前，应复制一份再修改，避免 Steam 更新覆盖配置。
4. 涉及参数、默认值或行为时，以官方 Dedicated Server Guide 为准。
5. 不应凭空创建 `server.cfg`、`valheim_server.cfg` 等官方不存在的配置文件。

## 官方配置入口

安装：

Steam → Library → Tools → Valheim Dedicated Server → Install。

服务器安装目录：

Steam → Valheim Dedicated Server → Manage → Browse local files。

Windows 启动脚本：

`start_headless_server.bat`

Linux 启动脚本：

`start_server.sh`

官方文档指出，服务器启动参数位于启动脚本中以 `start valheim_server...` 开头的命令行。

建议：

```text
start_headless_server.bat
        ↓
复制
        ↓
start_headless_server_custom.bat
        ↓
修改启动参数
        ↓
启动服务器
```

## 基础启动参数

### `-name`

服务器名称。

```text
-name "My Valheim Server"
```

用于设置服务器列表中显示的名称。

### `-port`

服务器端口。

```text
-port 2456
```

默认使用：

```text
2456-2457
```

Valheim 使用指定端口以及 `port + 1`。

使用 Steam Backend 时，需要在路由器上进行端口转发。

使用 Crossplay Backend 时，不需要为公网访问配置路由器端口转发。

### `-world`

世界名称。

```text
-world "Dedicated"
```

如果世界不存在，会创建新世界。

如果已经存在，则加载已有世界。

### `-password`

服务器密码。

```text
-password "Secret"
```

### `-savedir`

修改默认存档目录。

```text
-savedir "D:\ValheimData"
```

默认 Windows 路径：

```text
../%USERPROFILE%/AppData/LocalLow/IronGate/Valheim
```

默认 Linux 路径：

```text
~/.config/unity3d/IronGate/Valheim
```

世界文件和权限文件都位于存档目录。

### `-public`

控制服务器是否显示在服务器列表。

```text
-public 1
```

默认：

```text
1
```

隐藏服务器：

```text
-public 0
```

`-public 0` 适合局域网服务器或不希望出现在服务器浏览器中的服务器。

### `-logFile`

设置日志文件路径。

```text
-logFile "D:\ValheimData\server.log"
```

### `-saveinterval`

设置世界自动保存间隔，单位为秒。

```text
-saveinterval 1800
```

默认：

```text
1800
```

即 30 分钟。

## 自动备份

### `-backups`

设置自动备份数量。

```text
-backups 4
```

默认情况下：

- 第一个备份为 short backup
- 其余为 long backup

默认配置相当于：

```text
1 个 2 小时前的备份
3 个每隔 12 小时的备份
```

### `-backupshort`

设置 short backup 间隔。

```text
-backupshort 7200
```

默认：

```text
7200 秒
```

即 2 小时。

### `-backuplong`

设置 long backup 间隔。

```text
-backuplong 43200
```

默认：

```text
43200 秒
```

即 12 小时。

## Crossplay

### `-crossplay`

启用 Crossplay Backend。

```text
-crossplay
```

启用后：

- 使用 PlayFab Relay
- 不要求公网端口转发
- 不同平台玩家可以加入

未指定 `-crossplay` 时使用 Steam Backend。

Steam Backend：

```text
Steam 玩家
+
端口转发
```

Crossplay：

```text
多平台
+
Relay
+
无需公网端口转发
```

注意：Crossplay 服务器不能通过本地 IP 或 loopback IP 连接。

## 多服务器实例

### `-instanceid`

当同一台机器上运行多个服务器并使用相同端口时，可以设置唯一实例 ID。

```text
-instanceid "1"
```

例如：

```text
服务器 A
-instanceid "1"

服务器 B
-instanceid "2"
```

## 世界预设

### `-preset`

设置世界修饰器预设。

```text
-preset hard
```

有效值：

```text
Normal
Casual
Easy
Hard
Hardcore
Immersive
Hammer
```

注意：

设置 `-preset` 会覆盖之前的其他世界修饰器。

如果同时使用 `-preset` 和 `-modifier`，应让 `-modifier` 位于 `-preset` 后面。

## 世界修饰器

### `-modifier`

格式：

```text
-modifier <类型> <值>
```

例如：

```text
-modifier raids none
```

### Combat

```text
veryeasy
easy
hard
veryhard
```

示例：

```text
-modifier combat hard
```

### DeathPenalty

```text
casual
veryeasy
easy
hard
hardcore
```

示例：

```text
-modifier deathpenalty casual
```

### Resources

```text
muchless
less
more
muchmore
most
```

示例：

```text
-modifier resources more
```

### Raids

```text
none
muchless
less
more
muchmore
```

示例：

```text
-modifier raids none
```

### Portals

```text
casual
hard
veryhard
```

示例：

```text
-modifier portals casual
```

## 世界修饰器开关

### `-setkey`

设置世界修饰器 checkbox key。

有效值：

```text
nobuildcost
playerevents
passivemobs
nomap
```

示例：

```text
-setkey nomap
```

多个参数组合时，应保持官方参数语义，不自行扩展不存在的 key。

## 权限管理

权限文件位于默认存档目录，主要包括：

```text
adminlist.txt
bannedlist.txt
permittedlist.txt
```

作用：

```text
adminlist.txt
管理员

bannedlist.txt
封禁玩家

permittedlist.txt
允许玩家
```

每行写入一个 Platform User ID。

格式：

```text
[Platform]_[User ID]
```

Platform User ID 可以从服务器日志或游戏 F2 面板获取。

注意：

`permittedlist.txt` 一旦加入玩家，会使服务器变成白名单服务器。

即：

```text
permittedlist.txt 非空
        ↓
列表中的玩家允许进入
        ↓
其他玩家无法进入
```

## 控制台命令

游戏内按：

```text
F5
```

打开控制台。

常用命令：

```text
Kick PLAYERNAME
Ban PLAYERNAME
Unban PLAYERNAME
Banned
```

## 停止服务器

不要直接关闭服务器窗口。

推荐：

```text
CTRL+C
```

正常停止服务器。

直接点击窗口 X 可能导致服务器进程继续运行。

## Docker

官方 Dedicated Server 自带：

```text
docker_start_server.sh
```

通常与：

```text
start_server.sh
```

位于同一目录。

运行：

```bash
./docker_start_server.sh start_server.sh
```

首次运行会构建容器环境，可能需要数分钟。

之后再次启动会明显更快。

默认 Docker 数据卷：

```text
valheim_server_data
```

也可以修改：

```text
DOCKER_DATA_VOLUME="${HOME}/valheim_storage"
```

此时数据存储在：

```text
~/valheim_storage
```

其中包含：

```text
worlds
adminlist.txt
bannedlist.txt
permittedlist.txt
```

## Linux 依赖

Linux Dedicated Server 需要额外依赖：

```text
libatomic1
libpulse-dev
libpulse0
```

服务器要求：

```text
GLIBC_2.29
GLIBCXX_3.4.26
```

如果系统版本过低，可以考虑 Docker。

## 配置生成规则

如果需要为“一键启动 + 可视化服务器管理器”生成启动命令，应把 UI 配置映射为启动参数。

推荐数据结构：

```json
{
  "name": "My Valheim Server",
  "port": 2456,
  "world": "Dedicated",
  "password": "Secret",
  "public": 1,
  "savedir": "",
  "logFile": "",
  "saveInterval": 1800,
  "backups": 4,
  "backupShort": 7200,
  "backupLong": 43200,
  "crossplay": false,
  "instanceId": "",
  "preset": "Normal",
  "modifiers": {
    "combat": "",
    "deathPenalty": "",
    "resources": "",
    "raids": "",
    "portals": ""
  },
  "keys": []
}
```

然后生成：

```text
valheim_server.exe
-name "My Valheim Server"
-port 2456
-world "Dedicated"
-password "Secret"
-public 1
-saveinterval 1800
-backups 4
-backupshort 7200
-backuplong 43200
```

如果：

```text
crossplay = true
```

则追加：

```text
-crossplay
```

如果设置世界预设：

```text
-preset Hard
```

如果设置世界修饰器，则在 `-preset` 后追加：

```text
-modifier combat hard
-modifier resources more
-modifier raids none
```

如果设置 checkbox：

```text
-setkey nomap
```

## 一键启动应用的配置模型

推荐将配置分成四层：

```text
ServerConfig
├── 基础信息
│   ├── name
│   ├── port
│   ├── world
│   └── password
│
├── 网络
│   ├── public
│   └── crossplay
│
├── 存档与备份
│   ├── savedir
│   ├── saveInterval
│   ├── backups
│   ├── backupShort
│   └── backupLong
│
└── 世界规则
    ├── preset
    ├── modifiers
    └── keys
```

这样可以直接作为服务器管理 UI 的数据模型。

## Mod 管理边界

原版 Dedicated Server Guide 不负责 BepInEx、Thunderstore、r2modman 等 Mod 管理。

因此：

```text
原版服务器配置
≠
Mod 配置
```

不要把 Mod 参数混入官方 Dedicated Server 参数模型。

建议独立：

```text
ServerConfig
ModConfig
```

其中：

```text
ServerConfig
    ↓
启动参数

ModConfig
    ↓
BepInEx / Thunderstore / Mod 配置
```

## 事实依据

本 Skill 的核心参数和服务器行为依据 Valheim Dedicated Server 官方指南：

- Dedicated Server 安装方式
- Windows/Linux 启动脚本
- Steam Backend / Crossplay Backend
- 端口 2456-2457
- 存档目录
- 权限文件
- Docker
- 控制台命令
- 启动参数
- 自动备份
- 世界预设
- 世界修饰器

官方指南随 Dedicated Server 程序提供。

当本 Skill 与新版本官方文档冲突时，以新版本官方文档为准。

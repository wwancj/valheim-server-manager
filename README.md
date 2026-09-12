# Valheim Server Manager

A desktop application for managing Valheim Dedicated Server, BepInEx, and Thunderstore Mods.

## Tech Stack

- **Desktop:** Wails
- **Backend:** Go
- **Frontend:** React + TypeScript + Ant Design
- **Database:** SQLite

## Quick Start

```bash
# Install dependencies
go mod tidy
cd frontend && npm install

# Development
wails dev

# Build
wails build
```

## Project Structure

```
valheim-server-manager/
├── app/                  # Go backend
│   ├── services/
│   ├── models/
│   ├── repository/
│   ├── process/
│   ├── steamcmd/
│   ├── valheim/
│   ├── thunderstore/
│   ├── mods/
│   ├── backup/
│   ├── system/
│   └── config/
├── frontend/             # React + TS
│   └── src/
│       ├── pages/
│       ├── components/
│       ├── services/
│       ├── stores/
│       ├── types/
│       └── utils/
├── build/
├── docs/
├── scripts/
├── PLAN.md
├── CLAUDE.md
└── README.md
```

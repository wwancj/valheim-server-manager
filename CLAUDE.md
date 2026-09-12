# Valheim Server Manager Development Rules

## General

- Use Go for all system-level operations.
- Use React + TypeScript for UI.
- Keep business logic out of React components.
- Do not execute shell commands directly from the frontend.
- All filesystem, process, SteamCMD, Mod and server operations must go through Go services.
- Do not use localStorage for persistent application data.
- Use SQLite for persistent metadata.
- Use typed models between frontend and backend.
- Do not hardcode OS-specific paths.
- Do not hardcode shell commands in UI code.

## Architecture

```
Frontend
    ↓
Wails binding
    ↓
Go Service
    ↓
Repository / Process / Filesystem / External API
```

Each domain must have its own service:

- ServerService
- SteamCMDService
- ModService
- ThunderstoreService
- ConfigService
- BackupService
- LogService
- SystemService

## Tech Stack

- Desktop: Wails
- Backend: Go
- Frontend: React + TypeScript
- UI: Ant Design
- Database: SQLite

## Data Storage

```
AppData/
└── ValheimServerManager/
    ├── config/
    ├── servers/
    ├── cache/
    ├── logs/
    └── backups/
```

## Mod Safety

Never delete a Mod immediately when disabling it.
Prefer reversible operations.
Always check dependencies, version conflicts, file existence, installation state before modifying the server.

## Server Safety

Never modify world files while the server is running unless the operation is explicitly safe.
Before destructive operations: stop server → create backup → perform operation → validate result.

## Error Handling

Never swallow errors.
Every backend operation must return structured errors.
Frontend must display actionable error messages.

## UI

Do not put business logic inside JSX.
Avoid large components.
Extract reusable components.
Avoid duplicated button/action logic.
Use TypeScript types instead of `any` wherever possible.

## Testing

Every core service must have tests.
Priority:
1. Server process management
2. SteamCMD
3. Mod dependency resolution
4. Mod installation
5. Configuration
6. Backup/restore

## Implementation

Do not implement the entire project in one step.
Complete one phase at a time.
After each phase:
1. Build
2. Test
3. Fix errors
4. Update documentation
5. Commit

Do not proceed to the next phase if the current phase is broken.

package process

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"valheim-server-manager/app/events"
	"valheim-server-manager/app/config"
)

// Manager handles Valheim server process lifecycle.
type Manager struct {
	mu      sync.RWMutex
	cmd     *exec.Cmd
	running bool
	pid     int
	emitter events.EventEmitter
}

// NewManager creates a new process manager.
func NewManager() *Manager {
	return &Manager{emitter: &events.NoopEmitter{}}
}

// SetEventEmitter sets the event emitter for real-time notifications.
func (m *Manager) SetEventEmitter(emitter events.EventEmitter) {
	m.emitter = emitter
}

// CommandBuilder builds the Valheim server launch command from a full config.
type CommandBuilder struct {
	Config config.ServerConfig
}

// BuildArgs returns the command-line arguments from the full config.
func (b *CommandBuilder) BuildArgs() []string {
	cfg := b.Config
	args := []string{
		"-name", cfg.Server.Name,
		"-port", fmt.Sprintf("%d", cfg.Server.Port),
		"-world", cfg.Server.World,
		"-password", cfg.Server.Password,
	}
	if cfg.Server.Public {
		args = append(args, "-public", "1")
	} else {
		args = append(args, "-public", "0")
	}
	if cfg.Server.Crossplay {
		args = append(args, "-crossplay")
	}
	if cfg.Server.SaveDir != "" {
		args = append(args, "-savedir", cfg.Server.SaveDir)
	}
	if cfg.Server.LogFile != "" {
		args = append(args, "-logFile", cfg.Server.LogFile)
	}

	// Backup settings
	if cfg.Backup.SaveInterval > 0 {
		args = append(args, "-saveinterval", fmt.Sprintf("%d", cfg.Backup.SaveInterval))
	}
	if cfg.Backup.Count > 0 {
		args = append(args, "-backups", fmt.Sprintf("%d", cfg.Backup.Count))
	}
	if cfg.Backup.ShortInterval > 0 {
		args = append(args, "-backupshort", fmt.Sprintf("%d", cfg.Backup.ShortInterval))
	}
	if cfg.Backup.LongInterval > 0 {
		args = append(args, "-backuplong", fmt.Sprintf("%d", cfg.Backup.LongInterval))
	}

	// World preset
	if cfg.World.Preset != "" {
		args = append(args, "-preset", cfg.World.Preset)
	}

	// World modifiers
	modifiers := map[string]string{
		"combat":       cfg.World.Combat,
		"deathpenalty": cfg.World.DeathPenalty,
		"resources":    cfg.World.Resources,
		"raids":        cfg.World.Raids,
		"portals":      cfg.World.Portals,
	}
	for key, val := range modifiers {
		if val != "" && val != "normal" {
			args = append(args, "-modifier", key, val)
		}
	}

	// World keys (advanced toggles)
	keys := map[string]bool{
		"nobuildcost":  cfg.World.Keys.NoBuildCost,
		"playerevents": cfg.World.Keys.PlayerEvents,
		"passivemobs":  cfg.World.Keys.PassiveMobs,
		"nomap":        cfg.World.Keys.NoMap,
	}
	for key, enabled := range keys {
		if enabled {
			args = append(args, "-setkey", key)
		}
	}

	return args
}

// Start launches the Valheim server process.
func (m *Manager) Start(ctx context.Context, serverDir string, builder *CommandBuilder) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("server is already running")
	}

	// Find the server executable
	exe := "valheim_server.x86_64"
	if runtime.GOOS == "windows" {
		exe = "valheim_server.exe"
	}

	exePath := filepath.Join(serverDir, exe)
	if _, err := os.Stat(exePath); os.IsNotExist(err) {
		return fmt.Errorf("server executable not found: %s", exePath)
	}

	cmd := exec.Command(exePath, builder.BuildArgs()...)
	cmd.Dir = serverDir
	configureProcess(cmd)

	// Capture stdout/stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	m.cmd = cmd
	m.running = true
	m.pid = cmd.Process.Pid

	m.emitter.Emit("server:status", map[string]interface{}{
		"status": "started",
		"pid":    m.pid,
	})

	// Read output in background
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				m.emitter.Emit("log:new", map[string]interface{}{
					"line": string(buf[:n]),
				})
			}
			if err != nil {
				break
			}
		}
	}()

	// Wait for process to exit
	go func() {
		err := cmd.Wait()
		m.mu.Lock()
		m.running = false
		m.pid = 0
		m.mu.Unlock()

		status := "stopped"
		if err != nil {
			status = "error"
		}
		m.emitter.Emit("server:status", map[string]interface{}{
			"status": status,
			"error":  err,
		})
	}()

	return nil
}

// Stop stops the running server process.
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running || m.cmd == nil || m.cmd.Process == nil {
		return fmt.Errorf("server is not running")
	}

	if err := m.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	m.running = false
	m.pid = 0

	m.emitter.Emit("server:status", map[string]interface{}{
		"status": "stopped",
	})

	return nil
}

// IsRunning returns whether the server is currently running.
func (m *Manager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// GetPID returns the PID of the running server process.
func (m *Manager) GetPID() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.pid
}

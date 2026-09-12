package process

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Manager handles Valheim server process lifecycle.
type Manager struct {
	mu      sync.RWMutex
	cmd     *exec.Cmd
	running bool
	pid     int
}

// NewManager creates a new process manager.
func NewManager() *Manager {
	return &Manager{}
}

// CommandBuilder builds the Valheim server launch command.
type CommandBuilder struct {
	Name      string
	Port      int
	World     string
	Password  string
	Public    bool
	Preset    string
	Modifiers map[string]string
}

// BuildArgs returns the command-line arguments.
func (b *CommandBuilder) BuildArgs() []string {
	args := []string{
		"-name", b.Name,
		"-port", fmt.Sprintf("%d", b.Port),
		"-world", b.World,
		"-password", b.Password,
	}
	if b.Public {
		args = append(args, "-public", "1")
	} else {
		args = append(args, "-public", "0")
	}
	if b.Preset != "" {
		args = append(args, "-preset", b.Preset)
	}
	return args
}

// Start launches the Valheim server process.
func (m *Manager) Start(ctx context.Context, serverDir string, builder *CommandBuilder) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return fmt.Errorf("server is already running (PID: %d)", m.pid)
	}

	var exe string
	if runtime.GOOS == "windows" {
		exe = filepath.Join(serverDir, "valheim_server.exe")
	} else {
		exe = filepath.Join(serverDir, "valheim_server.x86_64")
	}

	if _, err := os.Stat(exe); os.IsNotExist(err) {
		return fmt.Errorf("server executable not found: %s", exe)
	}

	cmd := exec.Command(exe, builder.BuildArgs()...)
	cmd.Dir = serverDir

	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}

	m.cmd = cmd
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	m.running = true
	m.pid = cmd.Process.Pid

	wailsRuntime.EventsEmit(ctx, "server:status", map[string]interface{}{
		"status": "Running", "pid": m.pid,
	})

	go func() {
		err := cmd.Wait()
		m.mu.Lock()
		m.running = false
		m.pid = 0
		m.cmd = nil
		m.mu.Unlock()

		status := "Stopped"
		if err != nil {
			status = "Error"
		}
		wailsRuntime.EventsEmit(ctx, "server:status", map[string]interface{}{
			"status": status, "error": err,
		})
	}()

	return nil
}

// Stop gracefully stops the server.
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running || m.cmd == nil || m.cmd.Process == nil {
		return fmt.Errorf("server is not running")
	}

	wailsRuntime.EventsEmit(ctx, "server:status", map[string]interface{}{"status": "Stopping"})

	if err := m.cmd.Process.Signal(os.Interrupt); err != nil {
		if err := m.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}
	}
	return nil
}

// IsRunning returns whether the server process is running.
func (m *Manager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// GetPID returns the server process PID, or 0 if not running.
func (m *Manager) GetPID() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.pid
}

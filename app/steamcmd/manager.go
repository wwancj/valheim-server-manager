package steamcmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"valheim-server-manager/app/utils"
)

const (
	ValheimAppID     = "896660"
	steamcmdURLWin   = "https://steamcdn-a.akamaihd.net/client/installer/steamcmd.zip"
	steamcmdExeWin   = "steamcmd.exe"
	steamcmdExeLinux = "steamcmd.sh"
)

// Manager handles SteamCMD operations.
type Manager struct {
	installDir string
}

// NewManager creates a new SteamCMD manager.
func NewManager(installDir string) *Manager {
	if installDir == "" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, _ := os.UserHomeDir()
			appData = filepath.Join(home, ".config")
		}
		installDir = filepath.Join(appData, "ValheimServerManager", "steamcmd")
	}
	return &Manager{installDir: installDir}
}

// IsInstalled checks if SteamCMD is installed.
func (m *Manager) IsInstalled() bool {
	_, err := os.Stat(m.exePath())
	return err == nil
}

func (m *Manager) exePath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(m.installDir, steamcmdExeWin)
	}
	return filepath.Join(m.installDir, steamcmdExeLinux)
}

// Install downloads and installs SteamCMD.
func (m *Manager) Install(ctx context.Context) error {
	if m.IsInstalled() {
		return nil
	}
	if err := os.MkdirAll(m.installDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if runtime.GOOS == "windows" {
		return m.installWindows(ctx)
	}
	return m.installLinux(ctx)
}

func (m *Manager) installWindows(ctx context.Context) error {
	zipPath := filepath.Join(m.installDir, "steamcmd.zip")
	wailsRuntime.EventsEmit(ctx, "steamcmd:progress", map[string]interface{}{
		"stage": "downloading", "message": "正在下载 SteamCMD...", "progress": 0,
	})
	if err := utils.DownloadFile(steamcmdURLWin, zipPath); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	wailsRuntime.EventsEmit(ctx, "steamcmd:progress", map[string]interface{}{
		"stage": "extracting", "message": "正在解压...", "progress": 50,
	})
	if err := utils.Unzip(zipPath, m.installDir); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}
	os.Remove(zipPath)
	wailsRuntime.EventsEmit(ctx, "steamcmd:progress", map[string]interface{}{
		"stage": "done", "message": "SteamCMD 安装完成", "progress": 100,
	})
	return nil
}

func (m *Manager) installLinux(ctx context.Context) error {
	tarPath := filepath.Join(m.installDir, "steamcmd.tar.gz")
	wailsRuntime.EventsEmit(ctx, "steamcmd:progress", map[string]interface{}{
		"stage": "downloading", "message": "正在下载 SteamCMD...", "progress": 0,
	})
	if err := utils.DownloadFile("https://steamcdn-a.akamaihd.net/client/installer/steamcmd_linux.tar.gz", tarPath); err != nil {
		return err
	}
	wailsRuntime.EventsEmit(ctx, "steamcmd:progress", map[string]interface{}{
		"stage": "extracting", "message": "正在解压...", "progress": 50,
	})
	cmd := exec.Command("tar", "-xzf", tarPath, "-C", m.installDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}
	os.Remove(tarPath)
	wailsRuntime.EventsEmit(ctx, "steamcmd:progress", map[string]interface{}{
		"stage": "done", "message": "SteamCMD 安装完成", "progress": 100,
	})
	return nil
}

// InstallValheimServer downloads/updates the Valheim Dedicated Server.
func (m *Manager) InstallValheimServer(ctx context.Context, serverDir string) error {
	if !m.IsInstalled() {
		return fmt.Errorf("SteamCMD is not installed")
	}
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return err
	}

	wailsRuntime.EventsEmit(ctx, "steamcmd:valheim:progress", map[string]interface{}{
		"stage": "starting", "message": "正在启动 SteamCMD...", "progress": 0,
	})

	exe := m.exePath()
	args := []string{
		"+force_install_dir", serverDir,
		"+login", "anonymous",
		"+app_update", ValheimAppID, "validate",
		"+quit",
	}
	cmd := exec.Command(exe, args...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start SteamCMD: %w", err)
	}

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				line := string(buf[:n])
				wailsRuntime.EventsEmit(ctx, "steamcmd:valheim:output", line)
				if p := parseProgress(line); p >= 0 {
					wailsRuntime.EventsEmit(ctx, "steamcmd:valheim:progress", map[string]interface{}{
						"stage": "installing", "message": fmt.Sprintf("安装中... %d%%", p), "progress": p,
					})
				}
			}
			if err != nil {
				break
			}
		}
	}()
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				wailsRuntime.EventsEmit(ctx, "steamcmd:valheim:output", string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("SteamCMD exited with error: %w", err)
	}
	wailsRuntime.EventsEmit(ctx, "steamcmd:valheim:progress", map[string]interface{}{
		"stage": "done", "message": "Valheim Server 安装完成", "progress": 100,
	})
	return nil
}

func parseProgress(line string) int {
	if strings.Contains(line, "progress:") {
		idx := strings.Index(line, "progress:")
		var pct float64
		fmt.Sscanf(line[idx+9:], "%f", &pct)
		return int(pct)
	}
	return -1
}

package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"valheim-server-manager/app/utils"
)

const (
	bepinexDownloadURL = "https://github.com/BepInEx/BepInEx/releases/download/v5.4.2302/BepInEx_win_x64_5.4.2302.zip"
	bepinexDirName     = "BepInEx"
)

// BepInExService manages BepInEx installation.
type BepInExService struct{}

func NewBepInExService() *BepInExService {
	return &BepInExService{}
}

func (s *BepInExService) IsInstalled(serverDir string) bool {
	info, err := os.Stat(filepath.Join(serverDir, bepinexDirName))
	return err == nil && info.IsDir()
}

func (s *BepInExService) Install(ctx context.Context, serverDir string) error {
	if s.IsInstalled(serverDir) {
		wailsRuntime.EventsEmit(ctx, "bepinex:progress", map[string]interface{}{
			"stage": "done", "message": "BepInEx 已安装",
		})
		return nil
	}

	wailsRuntime.EventsEmit(ctx, "bepinex:progress", map[string]interface{}{
		"stage": "downloading", "message": "正在下载 BepInEx...", "progress": 0,
	})

	tmpZip := filepath.Join(os.TempDir(), "bepinex.zip")
	if err := utils.DownloadFile(bepinexDownloadURL, tmpZip); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer os.Remove(tmpZip)

	wailsRuntime.EventsEmit(ctx, "bepinex:progress", map[string]interface{}{
		"stage": "extracting", "message": "正在安装 BepInEx...", "progress": 50,
	})

	if err := utils.Unzip(tmpZip, serverDir); err != nil {
		return fmt.Errorf("extract failed: %w", err)
	}

	wailsRuntime.EventsEmit(ctx, "bepinex:progress", map[string]interface{}{
		"stage": "done", "message": "BepInEx 安装完成", "progress": 100,
	})
	return nil
}

func (s *BepInExService) GetPlugins(serverDir string) []string {
	pluginsDir := filepath.Join(serverDir, bepinexDirName, "plugins")
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		return []string{}
	}
	var plugins []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".dll") {
			plugins = append(plugins, e.Name())
		}
	}
	return plugins
}

package services

import (
	"fmt"
	"os"
	"path/filepath"
)

// SafetyCheck represents a single pre-launch check result.
type SafetyCheck struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // ok, warning, error
	Message string `json:"message"`
}

// SafetyService performs pre-launch safety checks.
type SafetyService struct{}

// NewSafetyService creates a new safety service.
func NewSafetyService() *SafetyService {
	return &SafetyService{}
}

// RunChecks performs all safety checks before server launch.
func (s *SafetyService) RunChecks(serverDir, password string, port int, bepInExInstalled bool, installedMods []string) []SafetyCheck {
	var checks []SafetyCheck

	// Check server directory
	checks = append(checks, s.checkServerDir(serverDir))

	// Check password
	checks = append(checks, s.checkPassword(password))

	// Check port
	checks = append(checks, s.checkPort(port))

	// Check world files
	checks = append(checks, s.checkWorldFiles(serverDir))

	// Check BepInEx
	checks = append(checks, s.checkBepInEx(serverDir, bepInExInstalled))

	return checks
}

func (s *SafetyService) checkServerDir(serverDir string) SafetyCheck {
	info, err := os.Stat(serverDir)
	if err != nil {
		return SafetyCheck{Name: "服务器目录", Status: "error", Message: "目录不存在: " + serverDir}
	}
	if !info.IsDir() {
		return SafetyCheck{Name: "服务器目录", Status: "error", Message: "路径不是目录"}
	}
	return SafetyCheck{Name: "服务器目录", Status: "ok", Message: "目录存在"}
}

func (s *SafetyService) checkPassword(password string) SafetyCheck {
	if password == "" {
		return SafetyCheck{Name: "服务器密码", Status: "warning", Message: "未设置密码，服务器将公开可访问"}
	}
	if len(password) < 5 {
		return SafetyCheck{Name: "服务器密码", Status: "warning", Message: "密码过短，建议至少5个字符"}
	}
	return SafetyCheck{Name: "服务器密码", Status: "ok", Message: "密码已设置"}
}

func (s *SafetyService) checkPort(port int) SafetyCheck {
	if port <= 0 || port > 65535 {
		return SafetyCheck{Name: "端口", Status: "error", Message: fmt.Sprintf("端口无效: %d", port)}
	}
	if port < 1024 {
		return SafetyCheck{Name: "端口", Status: "warning", Message: "使用系统端口(<1024)可能需要管理员权限"}
	}
	return SafetyCheck{Name: "端口", Status: "ok", Message: fmt.Sprintf("端口 %d", port)}
}

func (s *SafetyService) checkWorldFiles(serverDir string) SafetyCheck {
	savesDir := filepath.Join(serverDir, "saves")
	if _, err := os.Stat(savesDir); os.IsNotExist(err) {
		return SafetyCheck{Name: "世界文件", Status: "warning", Message: "saves 目录不存在，首次启动将自动创建"}
	}
	return SafetyCheck{Name: "世界文件", Status: "ok", Message: "世界文件目录存在"}
}

func (s *SafetyService) checkBepInEx(serverDir string, installed bool) SafetyCheck {
	if !installed {
		return SafetyCheck{Name: "BepInEx", Status: "ok", Message: "未安装（可选）"}
	}
	bepInExDir := filepath.Join(serverDir, "BepInEx")
	if _, err := os.Stat(bepInExDir); os.IsNotExist(err) {
		return SafetyCheck{Name: "BepInEx", Status: "warning", Message: "BepInEx 标记为已安装但目录不存在"}
	}
	return SafetyCheck{Name: "BepInEx", Status: "ok", Message: "已安装"}
}

// HasErrors checks if any check has an error status.
func HasErrors(checks []SafetyCheck) bool {
	for _, c := range checks {
		if c.Status == "error" {
			return true
		}
	}
	return false
}

// HasWarnings checks if any check has a warning status.
func HasWarnings(checks []SafetyCheck) bool {
	for _, c := range checks {
		if c.Status == "warning" {
			return true
		}
	}
	return false
}

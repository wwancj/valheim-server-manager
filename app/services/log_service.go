package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// LogService manages server log streaming.
type LogService struct {
	ctx      context.Context
	mu       sync.Mutex
	lines    []string
	maxLines int
	paused   bool
}

func NewLogService() *LogService {
	return &LogService{maxLines: 5000, lines: []string{}}
}

func (s *LogService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *LogService) AppendLog(line string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().Format("15:04:05")
	formatted := fmt.Sprintf("[%s] %s", timestamp, line)
	s.lines = append(s.lines, formatted)

	if len(s.lines) > s.maxLines {
		s.lines = s.lines[len(s.lines)-s.maxLines:]
	}

	if s.ctx != nil && !s.paused {
		wailsRuntime.EventsEmit(s.ctx, "log:new", formatted)
	}
}

func (s *LogService) GetLogs() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]string, len(s.lines))
	copy(result, s.lines)
	return result
}

func (s *LogService) ClearLogs() {
	s.mu.Lock()
	s.lines = []string{}
	s.mu.Unlock()
	if s.ctx != nil {
		wailsRuntime.EventsEmit(s.ctx, "log:clear", nil)
	}
}

func (s *LogService) SetPaused(paused bool) {
	s.mu.Lock()
	s.paused = paused
	s.mu.Unlock()
}

func (s *LogService) ExportLogs(dir string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if dir == "" {
		dir = filepath.Join(os.TempDir(), "valheim-logs")
	}
	os.MkdirAll(dir, 0755)

	filename := fmt.Sprintf("server_%s.log", time.Now().Format("20060102_150405"))
	path := filepath.Join(dir, filename)

	data := ""
	for _, line := range s.lines {
		data += line + "\n"
	}
	return path, os.WriteFile(path, []byte(data), 0644)
}

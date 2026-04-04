package log

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	TraceLevel = log.DebugLevel
	DebugLevel = log.DebugLevel
	InfoLevel  = log.InfoLevel
	WarnLevel  = log.WarnLevel
	ErrorLevel = log.ErrorLevel
)

var Logger *log.Logger

func init() {
	Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      "2006-01-02 15:04:05.000",
		Prefix:          "MoFrp-CLI",
		CallerOffset:    1,
	})
	styles := log.DefaultStyles()
	styles.Levels[TraceLevel] = lipgloss.NewStyle().
		SetString("TRACE").
		Bold(true).
		MaxWidth(5).
		Foreground(lipgloss.Color("61"))
	Logger.SetStyles(styles)
}

func InitLogger(logPath string, levelStr string, maxDays int, disableLogColor bool) {
	var output io.Writer
	var err error

	if logPath == "console" {
		output = os.Stdout
	} else {
		output, err = NewRotateFileWriter(logPath, maxDays)
		if err != nil {
			output = os.Stdout
		}
	}

	level, err := log.ParseLevel(levelStr)
	if err != nil {
		level = log.InfoLevel
	}

	Logger = log.NewWithOptions(output, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      "2006-01-02 15:04:05.000",
		Prefix:          "MoFrp-CLI",
		CallerOffset:    1,
		Level:           level,
	})
}

type RotateFileWriter struct {
	filePath    string
	maxDays     int
	file        *os.File
	lastRotate  time.Time
	currentDate string
}

func NewRotateFileWriter(filePath string, maxDays int) (*RotateFileWriter, error) {
	w := &RotateFileWriter{
		filePath:    filePath,
		maxDays:     maxDays,
		lastRotate:  time.Now(),
		currentDate: time.Now().Format("2006-01-02"),
	}

	if err := w.openFile(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *RotateFileWriter) openFile() error {
	var err error
	w.file, err = os.OpenFile(w.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	return err
}

func (w *RotateFileWriter) checkRotate() error {
	now := time.Now()
	currentDate := now.Format("2006-01-02")

	if currentDate != w.currentDate {
		if w.file != nil {
			w.file.Close()
		}

		oldPath := w.filePath
		newPath := w.filePath + "." + w.currentDate
		if _, err := os.Stat(oldPath); err == nil {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}

		w.cleanupOldLogs(now)

		w.currentDate = currentDate
		w.lastRotate = now
		return w.openFile()
	}

	return nil
}

func (w *RotateFileWriter) cleanupOldLogs(now time.Time) {
	if w.maxDays <= 0 {
		return
	}

	cutoffDate := now.AddDate(0, 0, -w.maxDays)

	dir := filepath.Dir(w.filePath)
	base := filepath.Base(w.filePath)

	files, _ := os.ReadDir(dir)
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		if dateStr, ok := strings.CutPrefix(name, base+"."); ok {
			if len(dateStr) == 10 {
				fileDate, err := time.Parse("2006-01-02", dateStr)
				if err == nil && fileDate.Before(cutoffDate) {
					os.Remove(filepath.Join(dir, name))
				}
			}
		}
	}
}

func (w *RotateFileWriter) Write(p []byte) (n int, err error) {
	if err := w.checkRotate(); err != nil {
		return 0, err
	}
	return w.file.Write(p)
}

func (w *RotateFileWriter) Close() error {
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

func Errorf(format string, v ...any) {
	Logger.Errorf(format, v...)
}

func Warnf(format string, v ...any) {
	Logger.Warnf(format, v...)
}

func Info(format string, v ...any) {
	Logger.Info(format, v...)
}

func Infof(format string, v ...any) {
	Logger.Infof(format, v...)
}

func Debugf(format string, v ...any) {
	Logger.Debugf(format, v...)
}

func Tracef(format string, v ...any) {
	Logger.Logf(TraceLevel, format, v...)
}

func Logf(level log.Level, offset int, format string, v ...any) {
	Logger.Logf(level, format, v...)
}

type WriteLogger struct {
	level  log.Level
	offset int
}

func NewWriteLogger(level log.Level, offset int) *WriteLogger {
	return &WriteLogger{
		level:  level,
		offset: offset,
	}
}

func (w *WriteLogger) Write(p []byte) (n int, err error) {
	msg := string(bytes.TrimRight(p, "\n"))
	Logger.Log(w.level, msg)
	return len(p), nil
}

func Fatalf(format string, v ...any) {
	Logger.Errorf(format, v...)
	os.Exit(1)
}

package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger zerolog.Logger

func InitLogger() {
	// File rotation
	rotatingFile := &lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    10, // MB
		MaxBackups: 5,
		MaxAge:     30, // days
		Compress:   true,
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	multi := io.MultiWriter(consoleWriter, rotatingFile)

	// 🔹 Customize caller path
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		projectName := "p2p-management-service" // change if needed

		rel := file
		if i := strings.Index(file, projectName); i != -1 {
			rel = file[i:] // keep everything from projectName onwards
		} else {
			rel = filepath.Base(file) // fallback: just filename
		}

		fn := runtime.FuncForPC(pc)
		if fn != nil {
			return fmt.Sprintf("%s:%d", rel, line)
		}
		return fmt.Sprintf("%s:%d", rel, line)
	}

	Logger = zerolog.New(multi).With().
		Timestamp().
		Caller().
		Logger()
}

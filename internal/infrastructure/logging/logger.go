package logging

import (
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

func Setup(serverMode string) {
	level := slog.LevelInfo
	if isDebugMode(serverMode) {
		level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)

	mode := normalizeGinMode(serverMode)
	gin.SetMode(mode)
	gin.DefaultWriter = os.Stdout
	gin.DefaultErrorWriter = os.Stderr
}

func GormLogger(serverMode string) logger.Interface {
	logLevel := logger.Error
	if isDebugMode(serverMode) {
		logLevel = logger.Info
	}

	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}

func isDebugMode(serverMode string) bool {
	return strings.EqualFold(serverMode, gin.DebugMode) ||
		strings.EqualFold(serverMode, "development")
}

func normalizeGinMode(serverMode string) string {
	switch strings.ToLower(serverMode) {
	case gin.DebugMode, "development":
		return gin.DebugMode
	case gin.TestMode:
		return gin.TestMode
	case gin.ReleaseMode, "":
		return gin.ReleaseMode
	default:
		return gin.ReleaseMode
	}
}

func Writer() io.Writer {
	return os.Stdout
}

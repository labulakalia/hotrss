package middle

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

// HTPPLog middle log use zero log
func HTPPLog(logfile string) gin.HandlerFunc {
	loglublogger := lumberjack.Logger{
		Filename: logfile,
	}
	httplog := zerolog.New(&loglublogger).With().Timestamp().Str("service", "hotrss").Logger()
	return logger.SetLogger(logger.Config{Logger: &httplog})
}

// RecoverLog middle log use zero log
func RecoverLog(logfile string) gin.HandlerFunc {
	loglublogger := lumberjack.Logger{
		Filename: logfile,
	}
	return gin.RecoveryWithWriter(&loglublogger)
}

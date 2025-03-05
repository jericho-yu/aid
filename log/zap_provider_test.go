package log

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test1(t *testing.T) {
	t.Run("test1 创建日志", func(t *testing.T) {
		logger := ZapProviderApp.New(
			ZapProviderConfig.New("./", false, EncoderTypeConsole, zapcore.ErrorLevel, true).
				SetMaxBackup(10).
				SetMaxDay(30).
				SetMaxSize(10),
		)

		logger.Info("test-info", zap.String("a", "b"))
		logger.Debug("test-debug", zap.String("c", "d"))
		logger.Warn("test-warning")
		logger.Error("test-error")
	})
}

package zapLogger

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test1(t *testing.T) {
	t.Run("test1 创建日志", func(t *testing.T) {
		zapLogger := ZapProviderApp.New(
			ZapProviderConfig.New("./", false, EncoderTypeConsole, zapcore.ErrorLevel).
				SetMaxBackup(30).
				SetMaxSize(10).
				SetMaxDay(30),
		)

		zapLogger.Info("test-info", zap.String("a", "b"))
		zapLogger.Debug("test-debug", zap.String("c", "d"))
		zapLogger.Warn("test-warning")
		zapLogger.Error("test-error")
	})
}

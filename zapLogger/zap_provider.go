package zapLogger

import (
	"fmt"
	"os"
	"time"

	"github.com/jericho-yu/aid/filesystem"
	"github.com/jericho-yu/aid/operation"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapProvider Zap日志服务提供者
type (
	ZapProvider struct {
		logger *zap.Logger
		err    error
	}

	EncoderType string
)

var ZapProviderApp ZapProvider

const (
	EncoderTypeConsole EncoderType = "CONSOLE"
	EncoderTypeJson    EncoderType = "JSON"
)

// getWriteSync 获取 zapcore.WriteSync
func getWriteSync(config *zapConfig, path string) zapcore.WriteSyncer {
	fileWriter := &lumberjack.Logger{
		Filename:   path,                // 日志文件名称
		MaxSize:    config.MaxSize,      // 文件大小限制,单位MB
		MaxBackups: config.MaxBackup,    // 最大保留日志文件数量
		MaxAge:     config.MaxDay,       // 日志文件保留天数
		Compress:   config.NeedCompress, // 是否压缩处理,压缩以后文件为xxxxx.gz
	}

	if config.InConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(fileWriter), zapcore.AddSync(os.Stdout))
	} else {
		return zapcore.AddSync(fileWriter)
	}
}

// New 实例化：Zap日志服务提供者
func (*ZapProvider) New(config *zapConfig) *ZapProvider { return NewZapProvider(config) }

// NewZapProvider 实例化：Zap日志服务提供者
//
//go:fix 推荐使用：New方法
func NewZapProvider(config *zapConfig) *ZapProvider {
	var (
		zapProvider     *ZapProvider = &ZapProvider{}
		fs              *filesystem.FileSystem
		zapCores        = make([]zapcore.Core, 0, 7)
		zapLoggerConfig = zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "time",
			NameKey:       "logger",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
				encoder.AppendString(t.Format(time.DateTime + ".000"))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		}
		encoderTypes = map[EncoderType]func(cfg zapcore.EncoderConfig) zapcore.Encoder{
			EncoderTypeJson:    zapcore.NewJSONEncoder,
			EncoderTypeConsole: zapcore.NewConsoleEncoder,
		}
	)

	fs = operation.TernaryFuncAll(func() bool { return config.PathAbs }, func() *filesystem.FileSystem { return filesystem.FileSystemApp.NewByAbs(config.Path) }, func() *filesystem.FileSystem { return filesystem.FileSystemApp.NewByRel(config.Path) })
	if !fs.IsExist {
		if zapProvider.err = fs.MkDir(); zapProvider.err != nil {
			zapProvider.err = fmt.Errorf("创建日志目录失败：%w", zapProvider.err)
			return zapProvider
		}
	}

	if config.Level < zapcore.DebugLevel {
		config.Level = zapcore.DebugLevel
	}

	if config.Level > zapcore.FatalLevel {
		config.Level = zapcore.FatalLevel
	}

	for _, logLevel := range []zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel} {
		if config.Level >= logLevel {
			writer := getWriteSync(config, fs.Copy().Join(logLevel.String()+config.Extension).GetDir())
			zapCores = append(zapCores, zapcore.NewCore(encoderTypes[config.EncoderType](zapLoggerConfig), writer, logLevel))
		}
	}

	zapProvider.logger = zap.New(zapcore.NewTee(zapCores...))

	defer func() {
		if config.InConsole {
			return
		}
		zapProvider.err = zapProvider.logger.Sync()
	}()

	return zapProvider
}

// Error 获取错误信息
func (my *ZapProvider) Error() error { return my.err }

// Logger 获取日志器
func (my *ZapProvider) Logger() *zap.Logger { return my.logger }

// func CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
// 	encoder.AppendString(t.Format(time.DateTime + ".000"))
// }

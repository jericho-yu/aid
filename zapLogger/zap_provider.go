package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jericho-yu/aid/filesystem"
	"github.com/jericho-yu/aid/operation"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapProvider Zap日志服务提供者
type (
	ZapProvider struct{}
	Cutter      struct {
		level    string        // 日志级别(debug, info, warn, error, dpanic, panic, fatal)
		format   string        // 时间格式(2006-01-02)
		Director string        // 日志文件夹
		file     *os.File      // 文件句柄
		mutex    *sync.RWMutex // 读写锁
	}

	zapConfig struct {
		Path         string
		PathAbs      bool
		MaxSize      int
		MaxBackup    int
		MaxDay       int
		NeedCompress bool
		InConsole    bool
		Extension    string
		Level        zapcore.Level
		EncoderType  EncoderType
	}

	EncoderType string
)

var (
	ZapProviderApp    ZapProvider
	ZapProviderConfig zapConfig
)

const (
	EncoderTypeConsole EncoderType = "CONSOLE"
	EncoderTypeJson    EncoderType = "JSON"
)

func newCutter(director string, level string, options ...func(*Cutter)) *Cutter {
	rotate := &Cutter{
		level:    level,
		Director: director,
		mutex:    new(sync.RWMutex),
	}
	for i := 0; i < len(options); i++ {
		options[i](rotate)
	}
	return rotate
}

// Write satisfies the io.Writer interface. It writes to the
// appropriate file handle that is currently being used.
// If we have reached rotation time, the target file gets
// automatically rotated, and also purged if necessary.
func (c *Cutter) Write(bytes []byte) (n int, err error) {
	c.mutex.Lock()
	defer func() {
		if c.file != nil {
			_ = c.file.Close()
			c.file = nil
		}
		c.mutex.Unlock()
	}()
	var business string
	if strings.Contains(string(bytes), "business") {
		var compile *regexp.Regexp
		compile, err = regexp.Compile(`{"business": "([^,]+)"}`)
		if err != nil {
			return 0, err
		}
		if compile.Match(bytes) {
			finds := compile.FindSubmatch(bytes)
			business = string(finds[len(finds)-1])
			bytes = compile.ReplaceAll(bytes, []byte(""))
		}
		compile, err = regexp.Compile(`"business": "([^,]+)"`)
		if err != nil {
			return 0, err
		}
		if compile.Match(bytes) {
			finds := compile.FindSubmatch(bytes)
			business = string(finds[len(finds)-1])
			bytes = compile.ReplaceAll(bytes, []byte(""))
		}
	}
	format := time.Now().Format(c.format)
	formats := make([]string, 0, 4)
	formats = append(formats, c.Director)
	if format != "" {
		formats = append(formats, format)
	}
	if business != "" {
		formats = append(formats, business)
	}
	formats = append(formats, c.level+".log")
	filename := filepath.Join(formats...)
	dirname := filepath.Dir(filename)
	err = os.MkdirAll(dirname, 0755)
	if err != nil {
		return 0, err
	}
	c.file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}

	return c.file.Write(bytes)
}

// GetWriteSync 获取 zapcore.WriteSync
func GetWriteSync(config *zapConfig, path string) zapcore.WriteSyncer {
	// fileWriter := newCutter(path, level.String(), func(c *Cutter) { c.format = time.DateOnly })

	fileWriter := &lumberjack.Logger{
		Filename:   path,                // 日志文件名称
		MaxSize:    config.MaxSize,      // 文件大小限制,单位MB
		MaxBackups: config.MaxBackup,    // 最大保留日志文件数量
		MaxAge:     config.MaxDay,       // 日志文件保留天数
		Compress:   config.NeedCompress, // 是否压缩处理,压缩以后文件为xxxxx.gz
	}

	return operation.Ternary[zapcore.WriteSyncer](config.InConsole, zapcore.NewMultiWriteSyncer(zapcore.AddSync(fileWriter), zapcore.AddSync(os.Stdout)), zapcore.AddSync(fileWriter))
}

// New 实例化：日志配置
func (*zapConfig) New(path string, pathAbs bool, encoderType EncoderType, level zapcore.Level, inConsole bool) *zapConfig {
	return &zapConfig{
		Path:         path,
		PathAbs:      pathAbs,
		EncoderType:  encoderType,
		Level:        level,
		MaxSize:      1,
		MaxBackup:    5,
		MaxDay:       30,
		NeedCompress: false,
		InConsole:    inConsole,
		Extension:    ".log",
	}
}

// SetMaxSize 设置单文件最大存储容量
func (my *zapConfig) SetMaxSize(maxSize int) *zapConfig {
	my.MaxSize = maxSize

	return my
}

// SetMaxBackup 设置最大备份数量
func (my *zapConfig) SetMaxBackup(maxBackup int) *zapConfig {
	my.MaxBackup = maxBackup

	return my
}

// SetMaxDay 设置日志文件最大保存天数
func (my *zapConfig) SetMaxDay(maxDay int) *zapConfig {
	my.MaxDay = maxDay

	return my
}

// SetNeedCompress 设置是否需要压缩
func (my *zapConfig) SetNeedCompress(needCompress bool) *zapConfig {
	my.NeedCompress = needCompress

	return my
}

func (*ZapProvider) New(config *zapConfig) *zap.Logger { return NewZapProvider(config) }

// NewZapProvider 实例化：Zap日志服务提供者
//
//go:fix 推荐使用：New方法
func NewZapProvider(config *zapConfig) *zap.Logger {
	var (
		e               error
		fs              *filesystem.FileSystem
		zapLogger       *zap.Logger
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

	if config.PathAbs {
		fs = filesystem.FileSystemApp.NewByAbs(config.Path)
	} else {
		fs = filesystem.FileSystemApp.NewByRel(config.Path)
	}
	fs = filesystem.FileSystemApp.NewByRelative(config.Path)
	if !fs.IsExist {
		e = fs.MkDir()
		if e != nil {
			panic(fmt.Errorf("创建日志目录失败：%s", e.Error()))
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
			writer := GetWriteSync(config, fs.Copy().Join(fmt.Sprintf("%s%s", logLevel.String(), config.Extension)).GetDir())
			zapCores = append(zapCores, zapcore.NewCore(encoderTypes[config.EncoderType](zapLoggerConfig), writer, logLevel))
		}
	}

	zapLogger = zap.New(zapcore.NewTee(zapCores...))

	defer func() {
		if config.InConsole {
			return
		}
		if e = zapLogger.Sync(); e != nil {
			panic(e)
		}
	}()

	return zapLogger
}

// func CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
// 	encoder.AppendString(t.Format(time.DateTime + ".000"))
// }

package log

import (
	"fmt"
	"github.com/jericho-yu/aid/filesystem"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ZapProvider Zap日志服务提供者
type (
	ZapProvider struct {
		path      string
		inConsole bool
	}
	fileRotateLogs struct{}
	Cutter         struct {
		level    string        // 日志级别(debug, info, warn, error, dpanic, panic, fatal)
		format   string        // 时间格式(2006-01-02)
		Director string        // 日志文件夹
		file     *os.File      // 文件句柄
		mutex    *sync.RWMutex // 读写锁
	}
	CutterOption func(*Cutter)

	ZapLoggerEncoderType string
)

const (
	ZapLoggerEncoderTypeConsole ZapLoggerEncoderType = "console"
	ZapLoggerEncoderTypeJson    ZapLoggerEncoderType = "json"
)

// WithCutterFormat 设置时间格式
func WithCutterFormat(format string) CutterOption {
	return func(c *Cutter) {
		c.format = format
	}
}

func NewCutter(director string, level string, options ...CutterOption) *Cutter {
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

var (
	zapProvider    *ZapProvider
	FileRotateLogs = new(fileRotateLogs)
)

// GetWriteSync 获取 zapcore.WriteSync
// Author [SliverHorn](https://github.com/SliverHorn)
func (r *fileRotateLogs) GetWriteSync(path, level string, inConsole bool) zapcore.WriteSyncer {
	fileWriter := NewCutter(path, level, WithCutterFormat(time.DateOnly))
	if inConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(fileWriter), zapcore.AddSync(os.Stdout))
	}
	return zapcore.AddSync(fileWriter)
}

// NewZapProvider 实例化：Zap日志服务提供者
func NewZapProvider(path string, inConsole bool, encoderType ZapLoggerEncoderType) *zap.Logger {
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
		zapLoggerEncoderTypes = map[ZapLoggerEncoderType]func(cfg zapcore.EncoderConfig) zapcore.Encoder{
			ZapLoggerEncoderTypeJson:    zapcore.NewJSONEncoder,
			ZapLoggerEncoderTypeConsole: zapcore.NewConsoleEncoder,
		}
	)

	if zapProvider == nil {
		zapProvider = &ZapProvider{path: path, inConsole: inConsole}
	}

	fs = filesystem.FileSystemApp.NewByRelative(path)
	if !fs.IsExist {
		e = fs.MkDir()
		if e != nil {
			panic(fmt.Errorf("创建日志目录失败：%s", e.Error()))
		}
	}

	for _, level := range []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.DPanicLevel,
		zapcore.PanicLevel,
		zapcore.FatalLevel,
	} {
		writer := FileRotateLogs.GetWriteSync(path, level.String(), inConsole)
		zapCores = append(zapCores, zapcore.NewCore(zapLoggerEncoderTypes[encoderType](zapLoggerConfig), writer, level))
	}

	zapLogger = zap.New(zapcore.NewTee(zapCores...))

	defer func() {
		if inConsole {
			return
		}
		e = zapLogger.Sync()
		if e != nil {
			panic(e)
		}
	}()

	return zapLogger
}

func CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format(time.DateTime + ".000"))
}

package zapLogger

import "go.uber.org/zap/zapcore"

type (
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

	zapConfigPath         struct{ Value string }
	zapConfigPathAbs      struct{ Value bool }
	zapConfigMaxSize      struct{ Value int }
	zapConfigMaxBackup    struct{ Value int }
	zapConfigMaxDay       struct{ Value int }
	zapConfigNeedCompress struct{ Value bool }
	zapConfigInConsole    struct{ Value bool }
	zapConfigExtension    struct{ Value string }

	ConfigType string
)

var ZapProviderConfig zapConfig

func NewZapConfigMaxSize(value int) *zapConfigMaxSize {
	return &zapConfigMaxSize{Value: value}
}

func NewZapConfigMaxBackup(value int) *zapConfigMaxBackup {
	return &zapConfigMaxBackup{Value: value}
}

func NewZapConfigMaxDay(value int) *zapConfigMaxDay {
	return &zapConfigMaxDay{Value: value}
}

func NewZapConfigNeedCompress(value bool) *zapConfigNeedCompress {
	return &zapConfigNeedCompress{Value: value}
}

func NewZapConfigInConsole(value bool) *zapConfigInConsole {
	return &zapConfigInConsole{Value: value}
}

func NewZapConfigExtension(value string) *zapConfigExtension {
	return &zapConfigExtension{Value: value}
}

// New 实例化：日志配置
func (*zapConfig) New(path string, pathAbs bool, encoderType EncoderType, level zapcore.Level, zapConfigItems ...any) *zapConfig {
	ins := &zapConfig{
		Path:         path,
		PathAbs:      pathAbs,
		EncoderType:  encoderType,
		Level:        level,
		MaxSize:      1,
		MaxBackup:    5,
		MaxDay:       30,
		NeedCompress: false,
		InConsole:    false,
		Extension:    ".log",
	}

	if len(zapConfigItems) > 0 {
		for _, item := range zapConfigItems {
			switch i := item.(type) {
			case *zapConfigMaxSize:
				ins.MaxSize = i.Value
			case *zapConfigMaxBackup:
				ins.MaxBackup = i.Value
			case *zapConfigMaxDay:
				ins.MaxDay = i.Value
			case *zapConfigExtension:
				ins.Extension = i.Value
			case *zapConfigNeedCompress:
				ins.NeedCompress = i.Value
			case *zapConfigInConsole:
				ins.InConsole = i.Value
			}
		}
	}

	return ins
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

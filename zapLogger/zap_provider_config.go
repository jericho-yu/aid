package zapLogger

import "go.uber.org/zap/zapcore"

type zapConfig struct {
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

var ZapProviderConfig zapConfig

// New 实例化：日志配置
func (*zapConfig) New(path string, pathAbs bool, encoderType EncoderType, level zapcore.Level) *zapConfig {
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

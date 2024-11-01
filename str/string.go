package str

import (
	"log"
	"os"
	"strings"

	"github.com/jericho-yu/aid/common"
)

type (
	Str struct {
		original string
	}

	TerminalLog struct {
		format string
		enable bool
	}
)

func NewStr(original string) *Str { return &Str{original: original} }

// PadLeftZeros 前置补零
func (my *Str) PadLeftZeros(length int) (string, error) {
	var (
		err error
		res strings.Builder = strings.Builder{}
	)

	if len(my.original) >= length {
		return my.original, nil
	}

	for i := 0; i < length-len(my.original); i++ {
		res.WriteRune('0')
	}

	if _, err = res.WriteString(my.original); err != nil {
		return "", err
	}

	return res.String(), nil
}

// PadRightZeros 后置补零
func (my *Str) PadRightZeros(length int) (string, error) {
	var (
		err error
		res strings.Builder = strings.Builder{}
	)

	if len(my.original) >= length {
		return my.original, nil
	}

	if _, err = res.WriteString(my.original); err != nil {
		return "", err
	}

	for i := 0; i < length-len(my.original); i++ {
		res.WriteRune('0')
	}

	return res.String(), nil
}

// PadRight 后置填充
func (my *Str) PadRight(length int, s string) string {
	my.original += strings.Repeat(s, length-(len(my.original)%length))
	return my.original
}

// PadLeft 前置补充
func (my *Str) PadLeft(length int, s string) string {
	my.original = strings.Repeat(s, length-(len(my.original)%length)) + s
	return my.original
}

// NewTerminalLog 实例化：控制台日志
func NewTerminalLog(format string) *TerminalLog {
	return &TerminalLog{format: format, enable: common.ToBool(os.Getenv("UTIL.STR.TERMINAL-LOG.ENABLE"))}
}

// Info 打印日志行
func (r *TerminalLog) Info(v ...any) {
	if !r.enable {
		return
	}

	log.Printf("：\n--------------------[INFO]--------------------\n"+r.format+"\n-----------------[INFO FINISH]-----------------\n", v...)
}

// Success 打印成功
func (r *TerminalLog) Success(v ...any) {
	if !r.enable {
		return
	}

	log.Printf("：\n--------------------[SUCCESS]--------------------\n"+r.format+"\n-----------------[SUCCESS FINISH]-----------------\n", v...)
}

// Wrong 打印错误
func (r *TerminalLog) Wrong(v ...any) {
	if !r.enable {
		return
	}

	log.Printf("：\n--------------------[WRONG]--------------------\n"+r.format+"\n-----------------[WRONG FINISH]-----------------\n", v...)
}

// Error 打印错误并终止程序
func (r *TerminalLog) Error(v ...any) {
	if !r.enable {
		return
	}

	log.Fatalf("：\n--------------------[ERROR]--------------------\n"+r.format+"\n-----------------[ERROR FINISH]-----------------\n", v...)
}

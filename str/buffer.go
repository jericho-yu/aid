package str

import (
	"bytes"
	"fmt"
)

type (
	Buffer        struct{ original *bytes.Buffer }
	BufferJoinAny struct {
		original []any
		sep      string
	}
	BufferJoinString struct {
		original []string
		sep      string
	}
	BufferJoinByte struct {
		original []byte
		sep      string
	}
	BufferJoinRune struct {
		original []rune
		sep      string
	}
)

var (
	BufferApp           Buffer
	BufferJoinAnyApp    BufferJoinAny
	BufferJoinStringApp BufferJoinString
	BufferJoinByteApp   BufferJoinByte
	BufferJoinRuneApp   BufferJoinRune
)

// NewByString 实例化：通过字符串
func (*Buffer) NewByString(original string) *Buffer { return &Buffer{bytes.NewBufferString(original)} }

// NewByBytes 实例化：通过字节码
func (*Buffer) NewByBytes(original []byte) *Buffer { return &Buffer{bytes.NewBuffer(original)} }

// JoinAny 追加任意到字符串，并使用分隔符
func (my *Buffer) JoinAny(value *BufferJoinAny) *Buffer {
	my.original.WriteString(value.ToString())

	return my
}

// JoinString 追加任意到字符串，并使用分隔符
func (my *Buffer) JoinString(value *BufferJoinString) *Buffer {
	my.original.WriteString(value.ToString())

	return my
}

// JoinByte 追加任意到字符串，并使用分隔符
func (my *Buffer) JoinByte(value *BufferJoinByte) *Buffer {
	my.original.WriteString(value.ToString())

	return my
}

// JoinRune 追加任意到字符串，并使用分隔符
func (my *Buffer) JoinRune(value *BufferJoinRune) *Buffer {
	my.original.WriteString(value.ToString())

	return my
}

// Any 追加任意内容到字符串
func (my *Buffer) Any(values ...any) *Buffer {
	for _, value := range values {
		fmt.Fprintf(my.original, "%s", value)
	}

	return my
}

// String 追加写入字符串
func (my *Buffer) String(stringList ...string) *Buffer {
	for _, s := range stringList {
		my.original.WriteString(s)
	}

	return my
}

// Byte 追加写入字节
func (my *Buffer) Byte(byteList ...byte) *Buffer {
	for _, b := range byteList {
		my.original.WriteByte(b)
	}

	return my
}

// Rune 追加写入字符
func (my *Buffer) Rune(runeList ...rune) *Buffer {
	for _, v := range runeList {
		my.original.WriteRune(v)
	}

	return my
}

// ToString 获取字符串
func (my *Buffer) ToString() string { return my.original.String() }

// ToBytes 获取字节码
func (my *Buffer) ToBytes() []byte { return my.original.Bytes() }

// ToPtr 获取字符串指针
func (my *Buffer) ToPtr() *string {
	ret := my.original.String()
	return &ret
}

func JoinAnyOption(values ...any) *BufferJoinAny { return BufferJoinAnyApp.New(values...) }

func (*BufferJoinAny) New(values ...any) *BufferJoinAny {
	return &BufferJoinAny{original: values}
}

func (my *BufferJoinAny) Sep(sep string) *BufferJoinAny {
	my.sep = sep
	return my
}

func (my *BufferJoinAny) ToString() string {
	var buffer bytes.Buffer
	for i, value := range my.original {
		if i > 0 {
			buffer.WriteString(my.sep)
		}
		fmt.Fprintf(&buffer, "%s", value)
	}
	return buffer.String()
}

func JoinStringOption(values ...string) *BufferJoinString { return BufferJoinStringApp.New(values...) }

func (*BufferJoinString) New(values ...string) *BufferJoinString {
	return &BufferJoinString{original: values}
}

func (my *BufferJoinString) Sep(sep string) *BufferJoinString {
	my.sep = sep
	return my
}

func (my *BufferJoinString) ToString() string {
	var buffer bytes.Buffer
	for i, value := range my.original {
		if i > 0 {
			buffer.WriteString(my.sep)
		}
		buffer.WriteString(value)
	}
	return buffer.String()
}

func JoinByteOption(values ...byte) *BufferJoinByte { return BufferJoinByteApp.New(values...) }

func (*BufferJoinByte) New(values ...byte) *BufferJoinByte {
	return &BufferJoinByte{original: values}
}

func (my *BufferJoinByte) Sep(sep string) *BufferJoinByte {
	my.sep = sep
	return my
}

func (my *BufferJoinByte) ToString() string {
	var buffer bytes.Buffer
	for i, value := range my.original {
		if i > 0 {
			buffer.WriteString(my.sep)
		}
		buffer.WriteByte(value)
	}
	return buffer.String()
}

func JoinRuneOption(values ...rune) *BufferJoinRune { return BufferJoinRuneApp.New(values...) }

func (*BufferJoinRune) New(values ...rune) *BufferJoinRune {
	return &BufferJoinRune{original: values}
}

func (my *BufferJoinRune) Sep(sep string) *BufferJoinRune {
	my.sep = sep
	return my
}

func (my *BufferJoinRune) ToString() string {
	var buffer bytes.Buffer
	for i, value := range my.original {
		if i > 0 {
			buffer.WriteString(my.sep)
		}
		buffer.WriteRune(value)
	}
	return buffer.String()
}

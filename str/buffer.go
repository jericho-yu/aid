package str

import "bytes"

type Buffer struct{ original *bytes.Buffer }

var BufferApp Buffer

// NewByString 实例化：通过字符串
func (*Buffer) NewByString(original string) *Buffer { return &Buffer{bytes.NewBufferString(original)} }

// NewByBytes 实例化：通过字节码
func (*Buffer) NewByBytes(original []byte) *Buffer { return &Buffer{bytes.NewBuffer(original)} }

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

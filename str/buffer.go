package str

import "bytes"

type (
	Buffer struct {
		original *bytes.Buffer
	}
)

func NewByString(original string) *Buffer { return &Buffer{original: bytes.NewBufferString(original)} }

func NewByBytes(original []byte) *Buffer { return &Buffer{original: bytes.NewBuffer(original)} }

func (my *Buffer) WriteString(stringList ...string) *Buffer {
	for _, s := range stringList {
		my.original.WriteString(s)
	}
	return my
}

func (my *Buffer) WriterByte(byteList ...byte) *Buffer {
	for _, b := range byteList {
		my.original.WriteByte(b)
	}
	return my
}

func (my *Buffer) WriteRune(runeList ...rune) *Buffer {
	for _, v := range runeList {
		my.original.WriteRune(v)
	}
	return my
}

func (my *Buffer) ToString() string { return my.original.String() }

func (my *Buffer) ToBytes() []byte { return my.original.Bytes() }

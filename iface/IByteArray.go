package iface

type IByteArray interface {
	BufferLen() int

	ByteSliceLen() int

	ByteSliceCap() int

	Move()

	WriteIndex() int

	ReadIndex() int

	Read(n int) []byte

	Write(bytes []byte)

	Bytes() []byte

	CheckForMove()

	ReadWithoutReturn(n int)
}

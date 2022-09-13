package znet

type ByteArray struct {
	buf        []byte
	readIndex  int
	writeIndex int
}

func NewByteArray(len int) ByteArray {
	return ByteArray{
		buf:        make([]byte, len),
		readIndex:  0,
		writeIndex: 0,
	}
}

func (ba *ByteArray) BufferLen() int {
	return ba.writeIndex - ba.readIndex
}

func (ba *ByteArray) ByteSliceLen() int {
	return len(ba.buf)
}

func (ba *ByteArray) ByteSliceCap() int {
	return cap(ba.buf)
}

func (ba *ByteArray) Move() {
	newReadIndex := 0
	newWriteIndex := ba.BufferLen()
	copy(ba.buf, ba.buf[ba.readIndex:ba.writeIndex])
	ba.readIndex = newReadIndex
	ba.writeIndex = newWriteIndex
}

func (ba *ByteArray) WriteIndex() int {
	return ba.writeIndex
}

func (ba *ByteArray) ReadIndex() int {
	return ba.readIndex
}

func (ba *ByteArray) Read(n int) []byte {
	retBytes := ba.buf[ba.readIndex : ba.readIndex+n]

	ba.readIndex += n
	ba.CheckForMove()
	return retBytes
}

func (ba *ByteArray) Write(bytes []byte) {
	if len(ba.buf)-ba.writeIndex > len(bytes) {
		copy(ba.buf[ba.writeIndex:], bytes)
	} else {
		copy(ba.buf[ba.writeIndex:], bytes)
		for i := len(ba.buf) - ba.writeIndex; i < len(bytes); i++ {
			ba.buf = append(ba.buf, bytes[i])
		}
	}
	ba.writeIndex += len(bytes)
}

func (ba *ByteArray) Bytes() []byte {
	return ba.buf
}

func (ba *ByteArray) CheckForMove() {
	if ba.readIndex > len(ba.buf)/2 {
		ba.Move()
	}
}

func (ba *ByteArray) ReadWithoutReturn(n int) {
	ba.readIndex += n
	ba.CheckForMove()
}

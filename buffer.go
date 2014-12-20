package mynet

type Buffer struct {
	b []byte
}

func (buffer *Buffer) Get() []byte {
	return buffer.b
}

func (buffer *Buffer) Len() int {
	return len(buffer.b)
}

type ReadBuffer struct {
	Buffer
}

func NewReadBuffer(size int) *ReadBuffer {
	rb := new(ReadBuffer)
	rb.b = make([]byte, size)
	return rb
}

func (r *ReadBuffer) Read(b []byte) (n int, err error) {
	//@todo
	return
}

func (r *ReadBuffer) Bytes(n int) []byte {
	return r.b[:n]
}

type WriteBuffer struct {
	Buffer
}

func NewWriteBuffer(size int) *WriteBuffer {
	wb := new(WriteBuffer)
	wb.b = make([]byte, size)
	return wb
}

func (w *WriteBuffer) Write(b []byte) (n int, err error) {
	return
}

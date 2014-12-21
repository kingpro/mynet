package mynet

import (
	"encoding/binary"
	"errors"
	"io"
)

type Protocol interface {
	Read(r io.Reader, buf []byte) (n int, err error)
	Write(w io.Writer, b []byte, buf []byte) error
}

type SimpleProtocol struct {
	header             []byte
	MaxReadPacketSize  int
	MaxWritePacketSize int
}

func NewSimpleProtocol(maxReadPacketSize int, maxWritePacketSize int) *SimpleProtocol {
	return &SimpleProtocol{
		header:             make([]byte, 4),
		MaxReadPacketSize:  maxReadPacketSize,
		MaxWritePacketSize: maxWritePacketSize,
	}
}

func (p *SimpleProtocol) Read(r io.Reader, buf []byte) (n int, err error) {
	if _, err = io.ReadFull(r, p.header); err != nil {
		return
	}

	n = int(binary.BigEndian.Uint32(p.header))
	if n > p.MaxReadPacketSize {
		return 0, errors.New("packet too large")
	}

	if cap(buf) < n {
		buf = make([]byte, n)
	} else {
		buf = buf[:n]
	}

	if _, err = io.ReadFull(r, buf); err != nil {
		return
	}

	return
}

func (p *SimpleProtocol) Write(w io.Writer, b []byte, buf []byte) error {
	n := len(b)
	if n == 0 {
		return errors.New("write data is nil")
	}

	if n > p.MaxWritePacketSize {
		return errors.New("send packet too large")
	}

	packn := n + 4
	if cap(buf) < packn {
		buf = make([]byte, packn)
	}

	binary.BigEndian.PutUint32(buf[0:4], uint32(n))
	copy(buf[4:packn], b)
	_, err := w.Write(buf[:packn])
	return err
}

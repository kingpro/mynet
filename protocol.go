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

func NewSimpleProtocol(maxReadPacketSize int, maxWritePacketSize int, header int) *SimpleProtocol {
	return &SimpleProtocol{
		header:             make([]byte, header),
		MaxReadPacketSize:  maxReadPacketSize,
		MaxWritePacketSize: maxWritePacketSize,
	}
}

func (p *SimpleProtocol) Read(r io.Reader, buf []byte) (n int, err error) {
	if _, err = io.ReadFull(r, p.header); err != nil {
		return
	}

	switch cap(p.header) {
	case 2:
		n = int(binary.BigEndian.Uint16(p.header))
	case 4:
		n = int(binary.BigEndian.Uint32(p.header))
	default:
		return 0, errors.New("header err")
	}

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

	header := cap(p.header)

	packn := n + header
	if cap(buf) < packn {
		buf = make([]byte, packn)
	}

	switch header {
	case 2:
		binary.BigEndian.PutUint16(buf[0:2], uint16(n))
	case 4:
		binary.BigEndian.PutUint32(buf[0:4], uint32(n))
	default:
		return errors.New("header err")
	}

	copy(buf[header:packn], b)
	_, err := w.Write(buf[:packn])
	return err
}

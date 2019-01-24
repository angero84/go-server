package kbuffer

import (
	"io"
)

type KByteBuffer struct {
	buf		[]byte
	len		int
	w		int
	r		int
}

func NewKByteBuffer(len uint32) *KByteBuffer {

	return &KByteBuffer{
		buf:		make([]byte, len),
		len:		int(len),
	}
}

func (m *KByteBuffer) Len() int				{ return m.w-m.r }
func (m *KByteBuffer) Bytes() []byte		{ return m.buf[m.r:m.w]}


func (m *KByteBuffer) Write(p []byte) (n int, err error) {

	plen := len(p)
	need := m.w + plen

	if need > m.len {

		size := int(0)
		if 0 >= m.len {
			size = plen*2
		} else {
			rate := need / m.len
			size = m.len*rate*2
		}
		m.resize(size)
	}

	n = copy(m.buf[m.w:], p)
	m.w += n
	return
}

func (m *KByteBuffer) Read(p []byte) (n int, err error) {

	if m.r >= m.w {
		return 0, io.EOF
	}
	n = copy(p, m.buf[m.r:])
	m.r += n

	return
}

func (m *KByteBuffer) Next(n int) []byte {

	if m.r + n > m.w {
		return nil
	}

	m.r += n
	return m.buf[m.r-n:m.r]
}

func (m *KByteBuffer) BytesAfter(n int) []byte {

	if n >= m.w {
		return nil
	}

	return m.buf[n:m.w]
}

func (m *KByteBuffer) SetBucket(n int) {

	if n > m.len {
		if 0 >= m.len {
			m.resize(n*2)
		} else {
			m.resize(m.len*2)
		}
	}

	m.w = n
	m.r = 0
}

func (m *KByteBuffer) resize(len int) {

	slice := make([]byte, len)
	n := copy(slice, m.buf[:m.w])
	m.buf = slice
	m.len = len
	m.w = n

	if m.r > len {
		m.r = len
	}

	return
}
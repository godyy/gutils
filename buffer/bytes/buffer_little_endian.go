package bytes

import (
	"encoding/binary"
	"io"
)

func (b *Buffer) ReadLitInt16() (int16, error) {
	n, err := b.ReadLitUint16()
	return int16(n), err
}

func (b *Buffer) WriteLitInt16(i int16) error {
	return b.WriteLitUint16(uint16(i))
}

func (b *Buffer) ReadLitUint16() (i uint16, err error) {
	l := b.Readable()
	if l == 0 {
		return 0, io.EOF
	}
	if l < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	i = binary.LittleEndian.Uint16(b.buf[b.off : b.off+2])
	b.off += 2
	return
}

func (b *Buffer) WriteLitUint16(i uint16) error {
	m, ok := b.tryGrowByReslice(2)
	if !ok {
		m = b.grow(2)
	}
	binary.LittleEndian.PutUint16(b.buf[m:m+2], i)
	return nil
}

func (b *Buffer) ReadLitInt32() (int32, error) {
	n, err := b.ReadLitUint32()
	return int32(n), err
}

func (b *Buffer) WriteLitInt32(i int32) error {
	return b.WriteLitUint32(uint32(i))
}

func (b *Buffer) ReadLitUint32() (i uint32, err error) {
	l := b.Readable()
	if l == 0 {
		return 0, io.EOF
	}
	if l < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	i = binary.LittleEndian.Uint32(b.buf[b.off : b.off+4])
	b.off += 4
	return
}

func (b *Buffer) WriteLitUint32(i uint32) error {
	m, ok := b.tryGrowByReslice(4)
	if !ok {
		m = b.grow(4)
	}
	binary.LittleEndian.PutUint32(b.buf[m:m+4], i)
	return nil
}

func (b *Buffer) ReadLitInt64() (int64, error) {
	n, err := b.ReadLitUint64()
	return int64(n), err
}

func (b *Buffer) WriteLitInt64(i int64) error {
	return b.WriteLitUint64(uint64(i))
}

func (b *Buffer) ReadLitUint64() (i uint64, err error) {
	l := b.Readable()
	if l == 0 {
		return 0, io.EOF
	}
	if l < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	i = binary.LittleEndian.Uint64(b.buf[b.off : b.off+8])
	b.off += 8
	return
}

func (b *Buffer) WriteLitUint64(i uint64) error {
	m, ok := b.tryGrowByReslice(8)
	if !ok {
		m = b.grow(8)
	}
	binary.LittleEndian.PutUint64(b.buf[m:m+8], i)
	return nil
}

func (b *Buffer) WriteLitFloat32(f float32) error {
	return b.WriteFloat32(f, binary.LittleEndian)
}

func (b *Buffer) ReadLitFloat32() (f float32, err error) {
	return b.ReadFloat32(binary.LittleEndian)
}

func (b *Buffer) WriteLitFloat64(f float64) error {
	return b.WriteFloat64(f, binary.LittleEndian)
}

func (b *Buffer) ReadLitFloat64() (f float64, err error) {
	return b.ReadFloat64(binary.LittleEndian)
}

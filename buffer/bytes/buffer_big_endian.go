package bytes

import (
	"encoding/binary"
	"io"
)

func (b *Buffer) ReadBigInt16() (int16, error) {
	n, err := b.ReadBigUint16()
	return int16(n), err
}

func (b *Buffer) WriteBigInt16(i int16) error {
	return b.WriteBigUint16(uint16(i))
}

func (b *Buffer) ReadBigUint16() (i uint16, err error) {
	l := b.Readable()
	if l == 0 {
		return 0, io.EOF
	}
	if l < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	i = binary.BigEndian.Uint16(b.buf[b.off : b.off+2])
	b.off += 2
	return
}

func (b *Buffer) WriteBigUint16(i uint16) error {
	m, ok := b.tryGrowByReslice(2)
	if !ok {
		m = b.grow(2)
	}
	binary.BigEndian.PutUint16(b.buf[m:m+2], i)
	return nil
}

func (b *Buffer) ReadBigInt32() (int32, error) {
	n, err := b.ReadBigUint32()
	return int32(n), err
}

func (b *Buffer) WriteBigInt32(i int32) error {
	return b.WriteLitUint32(uint32(i))
}

func (b *Buffer) ReadBigUint32() (i uint32, err error) {
	l := b.Readable()
	if l == 0 {
		return 0, io.EOF
	}
	if l < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	i = binary.BigEndian.Uint32(b.buf[b.off : b.off+4])
	b.off += 4
	return
}

func (b *Buffer) WriteBigUint32(i uint32) error {
	m, ok := b.tryGrowByReslice(4)
	if !ok {
		m = b.grow(4)
	}
	binary.BigEndian.PutUint32(b.buf[m:m+4], i)
	return nil
}

func (b *Buffer) ReadBigInt64() (int64, error) {
	n, err := b.ReadBigUint64()
	return int64(n), err
}

func (b *Buffer) WriteBigInt64(i int64) error {
	return b.WriteBigUint64(uint64(i))
}

func (b *Buffer) ReadBigUint64() (i uint64, err error) {
	l := b.Readable()
	if l == 0 {
		return 0, io.EOF
	}
	if l < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	i = binary.BigEndian.Uint64(b.buf[b.off : b.off+8])
	b.off += 8
	return
}

func (b *Buffer) WriteBigUint64(i uint64) error {
	m, ok := b.tryGrowByReslice(8)
	if !ok {
		m = b.grow(8)
	}
	binary.BigEndian.PutUint64(b.buf[m:m+8], i)
	return nil
}

func (b *Buffer) WriteBigFloat32(f float32) error {
	return b.WriteFloat32(f, binary.BigEndian)
}

func (b *Buffer) ReadBigFloat32() (f float32, err error) {
	return b.ReadFloat32(binary.BigEndian)
}

func (b *Buffer) WriteBigFloat64(f float64) error {
	return b.WriteFloat64(f, binary.BigEndian)
}

func (b *Buffer) ReadBigFloat64() (f float64, err error) {
	return b.ReadFloat64(binary.BigEndian)
}

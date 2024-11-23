package bytes

import (
	"encoding/binary"
	"github.com/godyy/gutils/buffer"
	"io"
)

func (b *FixedBuffer) ReadBigInt16() (int16, error) {
	n, err := b.ReadBigUint16()
	return int16(n), err
}

func (b *FixedBuffer) WriteBigInt16(i int16) error {
	return b.WriteBigUint16(uint16(i))
}

func (b *FixedBuffer) ReadBigUint16() (i uint16, err error) {
	l := b.Readable()
	if l == 0 {
		return 0, io.EOF
	}
	if l < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	i = binary.BigEndian.Uint16(b.buf[b.r : b.r+2])
	b.r += 2
	return
}

func (b *FixedBuffer) WriteBigUint16(i uint16) error {
	l := b.Writable()
	if l == 0 {
		return buffer.ErrBufferFull
	}
	if l < 2 {
		return buffer.ErrExceedBufferLimit
	}

	b.slideReadable()

	binary.BigEndian.PutUint16(b.buf[b.w:b.w+2], i)
	b.w += 2
	return nil
}

func (b *FixedBuffer) ReadBigInt32() (int32, error) {
	n, err := b.ReadBigUint32()
	return int32(n), err
}

func (b *FixedBuffer) WriteBigInt32(i int32) error {
	return b.WriteBigUint32(uint32(i))
}

func (b *FixedBuffer) ReadBigUint32() (i uint32, err error) {
	l := b.Readable()
	if l == 0 {
		return 0, io.EOF
	}
	if l < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	i = binary.BigEndian.Uint32(b.buf[b.r : b.r+4])
	b.r += 4
	return
}

func (b *FixedBuffer) WriteBigUint32(i uint32) error {
	l := b.Writable()
	if l == 0 {
		return buffer.ErrBufferFull
	}
	if l < 4 {
		return buffer.ErrExceedBufferLimit
	}

	b.slideReadable()

	binary.BigEndian.PutUint32(b.buf[b.w:b.w+4], i)
	b.w += 4
	return nil
}

func (b *FixedBuffer) ReadBigInt64() (int64, error) {
	n, err := b.ReadBigUint64()
	return int64(n), err
}

func (b *FixedBuffer) WriteBigInt64(i int64) error {
	return b.WriteBigUint64(uint64(i))
}

func (b *FixedBuffer) ReadBigUint64() (i uint64, err error) {
	l := b.Readable()
	if l == 0 {
		return 0, io.EOF
	}
	if l < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	i = binary.BigEndian.Uint64(b.buf[b.r : b.r+8])
	b.r += 8
	return
}

func (b *FixedBuffer) WriteBigUint64(i uint64) error {
	l := b.Writable()
	if l == 0 {
		return buffer.ErrBufferFull
	}
	if l < 8 {
		return buffer.ErrExceedBufferLimit
	}

	b.slideReadable()

	binary.BigEndian.PutUint64(b.buf[b.w:b.w+8], i)
	b.w += 8
	return nil
}

func (b *FixedBuffer) WriteBigFloat32(f float32) error {
	return b.WriteFloat32(f, binary.BigEndian)
}

func (b *FixedBuffer) ReadBigFloat32() (f float32, err error) {
	return b.ReadFloat32(binary.BigEndian)
}

func (b *FixedBuffer) WriteBigFloat64(f float64) error {
	return b.WriteFloat64(f, binary.BigEndian)
}

func (b *FixedBuffer) ReadBigFloat64() (f float64, err error) {
	return b.ReadFloat64(binary.BigEndian)
}

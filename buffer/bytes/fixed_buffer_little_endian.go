package bytes

import (
	"encoding/binary"
	"github.com/godyy/gutils/buffer"
	"io"
)

func (b *FixedBuffer) ReadLitInt16() (int16, error) {
	n, err := b.ReadLitUint16()
	return int16(n), err
}

func (b *FixedBuffer) WriteLitInt16(i int16) error {
	return b.WriteLitUint16(uint16(i))
}

func (b *FixedBuffer) ReadLitUint16() (i uint16, err error) {
	l := b.Readable()
	if l == 0 {
		return 0, io.EOF
	}
	if l < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	i = binary.LittleEndian.Uint16(b.buf[b.r : b.r+2])
	b.r += 2
	return
}

func (b *FixedBuffer) WriteLitUint16(i uint16) error {
	l := b.Writable()
	if l == 0 {
		return buffer.ErrBufferFull
	}
	if l < 2 {
		return buffer.ErrExceedBufferLimit
	}

	b.slideReadable()

	binary.LittleEndian.PutUint16(b.buf[b.w:b.w+2], i)
	b.w += 2
	return nil
}

func (b *FixedBuffer) ReadLitInt32() (int32, error) {
	n, err := b.ReadLitUint32()
	return int32(n), err
}

func (b *FixedBuffer) WriteLitInt32(i int32) error {
	return b.WriteLitUint32(uint32(i))
}

func (b *FixedBuffer) ReadLitUint32() (i uint32, err error) {
	l := b.Readable()
	if l == 0 {
		return 0, io.EOF
	}
	if l < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	i = binary.LittleEndian.Uint32(b.buf[b.r : b.r+4])
	b.r += 4
	return
}

func (b *FixedBuffer) WriteLitUint32(i uint32) error {
	l := b.Writable()
	if l == 0 {
		return buffer.ErrBufferFull
	}
	if l < 4 {
		return buffer.ErrExceedBufferLimit
	}

	b.slideReadable()

	binary.LittleEndian.PutUint32(b.buf[b.w:b.w+4], i)
	b.w += 4
	return nil
}

func (b *FixedBuffer) ReadLitInt64() (int64, error) {
	n, err := b.ReadLitUint64()
	return int64(n), err
}

func (b *FixedBuffer) WriteLitInt64(i int64) error {
	return b.WriteLitUint64(uint64(i))
}

func (b *FixedBuffer) ReadLitUint64() (i uint64, err error) {
	l := b.Readable()
	if l == 0 {
		return 0, io.EOF
	}
	if l < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	i = binary.LittleEndian.Uint64(b.buf[b.r : b.r+8])
	b.r += 8
	return
}

func (b *FixedBuffer) WriteLitUint64(i uint64) error {
	l := b.Writable()
	if l == 0 {
		return buffer.ErrBufferFull
	}
	if l < 8 {
		return buffer.ErrExceedBufferLimit
	}

	b.slideReadable()

	binary.LittleEndian.PutUint64(b.buf[b.w:b.w+8], i)
	b.w += 8
	return nil
}

func (b *FixedBuffer) WriteLitFloat32(f float32) error {
	return b.WriteFloat32(f, binary.LittleEndian)
}

func (b *FixedBuffer) ReadLitFloat32() (f float32, err error) {
	return b.ReadFloat32(binary.LittleEndian)
}

func (b *FixedBuffer) WriteLitFloat64(f float64) error {
	return b.WriteFloat64(f, binary.LittleEndian)
}

func (b *FixedBuffer) ReadLitFloat64() (f float64, err error) {
	return b.ReadFloat64(binary.LittleEndian)
}

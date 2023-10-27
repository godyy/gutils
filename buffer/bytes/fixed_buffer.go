package bytes

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/godyy/gutils/buffer"
	"github.com/pkg/errors"
)

// FixedBuffer 定长字节缓冲区
type FixedBuffer struct {
	buf  []byte
	r, w int
}

func NewFixedBuffer(size int) *FixedBuffer {
	if size <= 0 {
		panic("bytes.NewFixedBuffer: size <= 0")
	}
	return &FixedBuffer{
		buf: make([]byte, size),
		r:   0,
		w:   0,
	}
}

func (b *FixedBuffer) Size() int {
	return len(b.buf)
}

func (b *FixedBuffer) Cap() int {
	return len(b.buf)
}

func (b *FixedBuffer) Reset() {
	b.r = 0
	b.w = 0
}

func (b *FixedBuffer) Readable() int {
	return b.w - b.r
}

func (b *FixedBuffer) Writable() int {
	return len(b.buf) - b.w + b.r
}

func (b *FixedBuffer) Data() []byte {
	return b.buf[:]
}

func (b *FixedBuffer) UnreadData() []byte {
	return b.buf[b.r:b.w]
}

func (b *FixedBuffer) slideReadable() {
	if b.r > 0 {
		copy(b.buf, b.buf[b.r:b.w])
		b.w -= b.r
		b.r = 0
	}
}

func (b *FixedBuffer) ReadByte() (c byte, err error) {
	l := b.Readable()
	if l == 0 {
		return 0, io.EOF
	}

	c = b.buf[b.r]
	b.r += 1
	return
}

func (b *FixedBuffer) WriteByte(c byte) error {
	if b.Writable() < 1 {
		return buffer.ErrBufferFull
	}

	b.slideReadable()

	b.buf[b.w] = c
	b.w += 1
	return nil
}

func (b *FixedBuffer) ReadInt8() (int8, error) {
	n, err := b.ReadUint8()
	return int8(n), err
}

func (b *FixedBuffer) WriteInt8(i int8) error {
	return b.WriteUint8(uint8(i))
}

func (b *FixedBuffer) ReadUint8() (i uint8, err error) {
	return b.ReadByte()
}

func (b *FixedBuffer) WriteUint8(i uint8) error {
	return b.WriteByte(i)
}

func (b *FixedBuffer) ReadInt16() (int16, error) {
	n, err := b.ReadUint16()
	return int16(n), err
}

func (b *FixedBuffer) WriteInt16(i int16) error {
	return b.WriteUint16(uint16(i))
}

func (b *FixedBuffer) ReadUint16() (i uint16, err error) {
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

func (b *FixedBuffer) WriteUint16(i uint16) error {
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

func (b *FixedBuffer) ReadInt32() (int32, error) {
	n, err := b.ReadUint32()
	return int32(n), err
}

func (b *FixedBuffer) WriteInt32(i int32) error {
	return b.WriteUint32(uint32(i))
}

func (b *FixedBuffer) ReadUint32() (i uint32, err error) {
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

func (b *FixedBuffer) WriteUint32(i uint32) error {
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

func (b *FixedBuffer) ReadInt64() (int64, error) {
	n, err := b.ReadUint64()
	return int64(n), err
}

func (b *FixedBuffer) WriteInt64(i int64) error {
	return b.WriteUint64(uint64(i))
}

func (b *FixedBuffer) ReadUint64() (i uint64, err error) {
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

func (b *FixedBuffer) WriteUint64(i uint64) error {
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

func (b *FixedBuffer) ReadBool() (bool, error) {
	c, err := b.ReadByte()
	if err != nil {
		return false, err
	}
	return c == 1, nil
}

func (b *FixedBuffer) WriteBool(v bool) error {
	if v {
		return b.WriteByte(1)
	} else {
		return b.WriteByte(0)
	}
}

func (b *FixedBuffer) ReadVarint16() (int16, error) {
	i, n := binary.Varint(b.buf[b.r:])
	if n == 0 {
		return 0, io.EOF
	}
	if n < 0 || n > MaxVarintLen16 {
		return 0, ErrVarintOverflow
	}
	b.r += n
	return int16(i), nil
}

func (b *FixedBuffer) WriteVarint16(i int16) error {
	var buf [MaxVarintLen16]byte
	n := binary.PutVarint(buf[:], int64(i))
	l := b.Writable()
	if l < n {
		return buffer.ErrExceedBufferLimit
	}
	b.slideReadable()
	copy(b.buf[b.w:], buf[:n])
	b.w += n
	return nil
}

func (b *FixedBuffer) ReadUvarint16() (uint16, error) {
	i, n := binary.Uvarint(b.buf[b.r:])
	if n == 0 {
		return 0, io.EOF
	}
	if n < 0 || n > MaxVarintLen16 {
		return 0, ErrVarintOverflow
	}
	b.r += n
	return uint16(i), nil
}

func (b *FixedBuffer) WriteUvarint16(i uint16) error {
	var buf [MaxVarintLen16]byte
	n := binary.PutUvarint(buf[:], uint64(i))
	l := b.Writable()
	if l < n {
		return buffer.ErrExceedBufferLimit
	}
	b.slideReadable()
	copy(b.buf[b.w:], buf[:n])
	b.w += n
	return nil
}

func (b *FixedBuffer) ReadVarint32() (int32, error) {
	i, n := binary.Varint(b.buf[b.r:])
	if n == 0 {
		return 0, io.EOF
	}
	if n < 0 || n > MaxVarintLen32 {
		return 0, ErrVarintOverflow
	}
	b.r += n
	return int32(i), nil
}

func (b *FixedBuffer) WriteVarint32(i int32) error {
	var buf [MaxVarintLen32]byte
	n := binary.PutVarint(buf[:], int64(i))
	l := b.Writable()
	if l < n {
		return buffer.ErrExceedBufferLimit
	}
	b.slideReadable()
	copy(b.buf[b.w:], buf[:n])
	b.w += n
	return nil
}

func (b *FixedBuffer) ReadUvarint32() (uint32, error) {
	i, n := binary.Uvarint(b.buf[b.r:])
	if n == 0 {
		return 0, io.EOF
	}
	if n < 0 || n > MaxVarintLen32 {
		return 0, ErrVarintOverflow
	}
	b.r += n
	return uint32(i), nil
}

func (b *FixedBuffer) WriteUvarint32(i uint32) error {
	var buf [MaxVarintLen32]byte
	n := binary.PutUvarint(buf[:], uint64(i))
	l := b.Writable()
	if l < n {
		return buffer.ErrExceedBufferLimit
	}
	b.slideReadable()
	copy(b.buf[b.w:], buf[:n])
	b.w += n
	return nil
}

func (b *FixedBuffer) ReadVarint64() (int64, error) {
	i, n := binary.Varint(b.buf[b.r:])
	if n == 0 {
		return 0, io.EOF
	}
	if n < 0 {
		return 0, ErrVarintOverflow
	}
	b.r += n
	return i, nil
}

func (b *FixedBuffer) WriteVarint64(i int64) error {
	var buf [MaxVarintLen64]byte
	n := binary.PutVarint(buf[:], i)
	l := b.Writable()
	if l < n {
		return buffer.ErrExceedBufferLimit
	}
	b.slideReadable()
	copy(b.buf[b.w:], buf[:n])
	b.w += n
	return nil
}

func (b *FixedBuffer) ReadUvarint64() (uint64, error) {
	i, n := binary.Uvarint(b.buf[b.r:])
	if n == 0 {
		return 0, io.EOF
	}
	if n < 0 {
		return 0, ErrVarintOverflow
	}
	b.r += n
	return i, nil
}

func (b *FixedBuffer) WriteUvarint64(i uint64) error {
	var buf [MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], i)
	l := b.Writable()
	if l < n {
		return buffer.ErrExceedBufferLimit
	}
	b.slideReadable()
	copy(b.buf[b.w:], buf[:n])
	b.w += n
	return nil
}

func (b *FixedBuffer) Read(p []byte) (int, error) {
	if len(p) <= 0 {
		return 0, io.ErrShortBuffer
	}

	if b.Readable() == 0 {
		return 0, io.EOF
	}

	n := copy(p, b.buf[b.r:b.w])
	b.r += n
	return n, nil
}

func (b *FixedBuffer) Write(p []byte) (n int, err error) {
	if len(p) <= 0 {
		return 0, nil
	}

	l := b.Writable()
	if l == 0 {
		return 0, buffer.ErrBufferFull
	}

	b.slideReadable()

	n = copy(b.buf[b.w:], p)
	b.w += n
	return
}

func (b *FixedBuffer) ReadString() (string, error) {
	i, n := binary.Varint(b.buf[b.r:])
	if n == 0 {
		return "", io.EOF
	}
	if n < 0 || n > MaxStringLenLen {
		return "", errors.WithMessage(ErrVarintOverflow, "read length")
	}

	l := int(i)
	if l < 0 {
		return "", buffer.ErrStringLenExceedLimit
	}

	if l+n > b.Readable() {
		return "", io.ErrUnexpectedEOF
	}

	b.r += n
	if l == 0 {
		return "", nil
	}
	s := string(b.buf[b.r : b.r+l])
	b.r += l
	return s, nil
}

func (b *FixedBuffer) WriteString(s string) error {
	l := len(s)
	if l > math.MaxInt32 {
		return buffer.ErrStringLenExceedLimit
	}

	var buf [MaxStringLenLen]byte
	ll := binary.PutVarint(buf[:], int64(l))

	if l+ll > b.Writable() {
		return buffer.ErrExceedBufferLimit
	}

	b.slideReadable()
	copy(b.buf[b.w:], buf[:ll])
	b.w += ll
	copy(b.buf[b.w:], s)
	b.w += l
	return nil
}

func (b *FixedBuffer) ReadFrom(r io.Reader) (int64, error) {
	if b.Writable() == 0 {
		return 0, buffer.ErrBufferFull
	}

	b.slideReadable()

	n, err := r.Read(b.buf[b.w:])
	if n < 0 {
		panic("bytes.FixedBuffer.ReadFrom: reader returned negative count from Read")
	}

	b.w += n
	return int64(n), err
}

func (b *FixedBuffer) WriteTo(w io.Writer) (int64, error) {
	l := b.Readable()
	if l == 0 {
		return 0, nil
	}

	n, err := w.Write(b.buf[b.r:b.w])
	if n > l {
		panic("bytes.FixedBuffer.WriteTo: invalid Write count")
	}

	b.r += n
	return int64(n), err
}

func (b *FixedBuffer) Peek(n int) ([]byte, error) {
	if n > len(b.buf) {
		return nil, buffer.ErrExceedBufferLimit
	}

	if n > b.Readable() {
		return b.buf[b.r:b.w], io.ErrUnexpectedEOF
	}

	return b.buf[b.r : b.r+n], nil
}

func (b *FixedBuffer) Skip(n int) (skipped int, err error) {
	if n > len(b.buf) {
		return 0, buffer.ErrExceedBufferLimit
	}

	if l := b.Readable(); n > l {
		skipped = l
		err = io.ErrUnexpectedEOF
	} else {
		skipped = n
	}

	b.r += skipped
	return
}

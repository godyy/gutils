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

// NewFixedBuffer 使用指定size创建FixedBuffer
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

// Size 获取buf的大小
func (b *FixedBuffer) Size() int {
	return len(b.buf)
}

// Reset 重置读写状态
func (b *FixedBuffer) Reset() {
	b.r = 0
	b.w = 0
}

// SetBuf 设置buf并返回之前的buf
func (b *FixedBuffer) SetBuf(buf []byte) []byte {
	old := b.buf
	b.buf = buf
	b.r = 0
	b.w = 0
	return old
}

// Readable 获取可读取数据长度
func (b *FixedBuffer) Readable() int {
	return b.w - b.r
}

// Writable 获取可写入数据长度
func (b *FixedBuffer) Writable() int {
	return len(b.buf) - b.w + b.r
}

// Data 获取buf中的完整写入数据
func (b *FixedBuffer) Data() []byte {
	if b.buf == nil {
		return nil
	}
	return b.buf[:b.w]
}

// UnreadData 获取buf中的未读数据
func (b *FixedBuffer) UnreadData() []byte {
	if b.buf == nil {
		return nil
	}
	return b.buf[b.r:b.w]
}

// slideReadable 将可读数据滑动到buf最前端
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

	i = binary.NativeEndian.Uint16(b.buf[b.r : b.r+2])
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

	binary.NativeEndian.PutUint16(b.buf[b.w:b.w+2], i)
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

	i = binary.NativeEndian.Uint32(b.buf[b.r : b.r+4])
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

	binary.NativeEndian.PutUint32(b.buf[b.w:b.w+4], i)
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

	i = binary.NativeEndian.Uint64(b.buf[b.r : b.r+8])
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

	binary.NativeEndian.PutUint64(b.buf[b.w:b.w+8], i)
	b.w += 8
	return nil
}

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

func (b *FixedBuffer) WriteVarint16(i int16) (int, error) {
	var buf [MaxVarintLen16]byte
	n := binary.PutVarint(buf[:], int64(i))
	l := b.Writable()
	if l < n {
		return 0, buffer.ErrExceedBufferLimit
	}
	b.slideReadable()
	copy(b.buf[b.w:], buf[:n])
	b.w += n
	return n, nil
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

func (b *FixedBuffer) WriteUvarint16(i uint16) (int, error) {
	var buf [MaxVarintLen16]byte
	n := binary.PutUvarint(buf[:], uint64(i))
	l := b.Writable()
	if l < n {
		return 0, buffer.ErrExceedBufferLimit
	}
	b.slideReadable()
	copy(b.buf[b.w:], buf[:n])
	b.w += n
	return n, nil
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

func (b *FixedBuffer) WriteVarint32(i int32) (int, error) {
	var buf [MaxVarintLen32]byte
	n := binary.PutVarint(buf[:], int64(i))
	l := b.Writable()
	if l < n {
		return 0, buffer.ErrExceedBufferLimit
	}
	b.slideReadable()
	copy(b.buf[b.w:], buf[:n])
	b.w += n
	return n, nil
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

func (b *FixedBuffer) WriteUvarint32(i uint32) (int, error) {
	var buf [MaxVarintLen32]byte
	n := binary.PutUvarint(buf[:], uint64(i))
	l := b.Writable()
	if l < n {
		return 0, buffer.ErrExceedBufferLimit
	}
	b.slideReadable()
	copy(b.buf[b.w:], buf[:n])
	b.w += n
	return n, nil
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

func (b *FixedBuffer) WriteVarint64(i int64) (int, error) {
	var buf [MaxVarintLen64]byte
	n := binary.PutVarint(buf[:], i)
	l := b.Writable()
	if l < n {
		return 0, buffer.ErrExceedBufferLimit
	}
	b.slideReadable()
	copy(b.buf[b.w:], buf[:n])
	b.w += n
	return n, nil
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

func (b *FixedBuffer) WriteUvarint64(i uint64) (int, error) {
	var buf [MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], i)
	l := b.Writable()
	if l < n {
		return 0, buffer.ErrExceedBufferLimit
	}
	b.slideReadable()
	copy(b.buf[b.w:], buf[:n])
	b.w += n
	return n, nil
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

// Peek 自buf中提取n个字节，同时并不会更新已读偏移
func (b *FixedBuffer) Peek(n int) ([]byte, error) {
	if n > len(b.buf) {
		return nil, buffer.ErrExceedBufferLimit
	}

	if n > b.Readable() {
		return b.buf[b.r:b.w], io.ErrUnexpectedEOF
	}

	return b.buf[b.r : b.r+n], nil
}

// Skip 自buf中跳过n个字节
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

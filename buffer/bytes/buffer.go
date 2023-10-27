package bytes

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/godyy/gutils/buffer"
	pkg_errors "github.com/pkg/errors"
)

// ErrBufferTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.
var ErrBufferTooLarge = errors.New("bytes: buffer too large")

// smallBufferSize is an initial allocation minimal capacity.
const smallBufferSize = 64

const maxInt = int(^uint(0) >> 1)

// Buffer 变长byte缓冲区
type Buffer struct {
	buf []byte
	off int
}

func NewBuffer(cap ...int) *Buffer {
	c := smallBufferSize
	if len(cap) > 0 {
		c = cap[0]
		if c <= 0 {
			panic("bytes.NewBuffer: cap <= 0")
		}
	}
	return &Buffer{
		buf: make([]byte, 0, c),
		off: 0,
	}
}

func NewBufferWithBuf(buf ...[]byte) *Buffer {
	b := &Buffer{
		buf: nil,
		off: 0,
	}
	if len(buf) > 0 {
		b.buf = buf[0]
	}
	if b.buf == nil {
		b.buf = make([]byte, 0, smallBufferSize)
	}
	return b
}

func (b *Buffer) Size() int {
	return len(b.buf)
}

func (b *Buffer) Cap() int {
	return cap(b.buf)
}

func (b *Buffer) Reset() {
	b.buf = b.buf[:0]
	b.off = 0
}

func (b *Buffer) ResetCap(minCap, maxCap int) {
	if minCap > maxCap {
		panic("bytes.Buffer.ResetCap: minCap > maxCap")
	}
	if minCap < smallBufferSize {
		minCap = smallBufferSize
	}

	if cap(b.buf) > maxCap {
		b.buf = make([]byte, 0, minCap)
		b.off = 0
	} else {
		b.Reset()
	}
}

func (b *Buffer) Readable() int {
	return len(b.buf) - b.off
}

func (b *Buffer) Writable() int {
	return cap(b.buf) - len(b.buf)
}

func (b *Buffer) Data() []byte {
	return b.buf[:]
}

func (b *Buffer) UnreadData() []byte {
	return b.buf[b.off:]
}

func (b *Buffer) tryGrowByReslice(n int) (int, bool) {
	if l := len(b.buf); n <= cap(b.buf)-l {
		b.buf = b.buf[:l+n]
		return l, true
	}
	return 0, false
}

// growBuffSlice grows b by n, preserving the original content of b.
// If the allocation fails, it panics with ErrTooLarge.
func growBuffSlice(b []byte, n int) []byte {
	defer func() {
		if recover() != nil {
			panic(ErrBufferTooLarge)
		}
	}()
	// TODO(http://golang.org/issue/51462): We should rely on the append-make
	// pattern so that the compiler can call runtime.growslice. For example:
	//	return append(b, make([]byte, n)...)
	// This avoids unnecessary zero-ing of the first len(b) bytes of the
	// allocated slice, but this pattern causes b to escape onto the heap.
	//
	// Instead use the append-make pattern with a nil slice to ensure that
	// we allocate buffers rounded up to the closest size class.
	c := len(b) + n // ensure enough space for n elements
	if c < 2*cap(b) {
		// The growth rate has historically always been 2x. In the future,
		// we could rely purely on append to determine the growth rate.
		c = 2 * cap(b)
	}
	b2 := append([]byte(nil), make([]byte, c)...)
	copy(b2, b)
	return b2[:len(b)]
}

func (b *Buffer) grow(n int) int {
	m := b.Readable()
	// If buffer is empty, reset to recover space.
	if m == 0 && b.off != 0 {
		b.Reset()
	}
	// Try to grow by means of a reslice.
	if i, ok := b.tryGrowByReslice(n); ok {
		return i
	}
	if b.buf == nil && n <= smallBufferSize {
		b.buf = make([]byte, n, smallBufferSize)
		return 0
	}
	c := cap(b.buf)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(b.buf, b.buf[b.off:])
	} else if c > maxInt-c-n {
		panic(ErrBufferTooLarge)
	} else {
		// Add b.off to account for b.buf[:b.off] being sliced off the front.
		b.buf = growBuffSlice(b.buf[b.off:], b.off+n)
	}
	// Restore b.off and len(b.buf).
	b.off = 0
	b.buf = b.buf[:m+n]

	return m
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the
// buffer without another allocation.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *Buffer) Grow(n int) {
	if n < 0 {
		panic("bytes.Buffer.Grow: n < 0")
	}
	m := b.grow(n)
	b.buf = b.buf[:m]
}

func (b *Buffer) ReadByte() (c byte, err error) {
	if b.Readable() == 0 {
		return 0, io.EOF
	}

	c = b.buf[b.off]
	b.off++
	return
}

func (b *Buffer) WriteByte(c byte) error {
	m, ok := b.tryGrowByReslice(1)
	if !ok {
		m = b.grow(1)
	}
	b.buf[m] = c
	return nil
}

func (b *Buffer) ReadInt8() (int8, error) {
	n, err := b.ReadUint8()
	return int8(n), err
}

func (b *Buffer) WriteInt8(i int8) error {
	return b.WriteUint8(uint8(i))
}

func (b *Buffer) ReadUint8() (i uint8, err error) {
	return b.ReadByte()
}

func (b *Buffer) WriteUint8(i uint8) error {
	return b.WriteByte(i)
}

func (b *Buffer) ReadInt16() (int16, error) {
	n, err := b.ReadUint16()
	return int16(n), err
}

func (b *Buffer) WriteInt16(i int16) error {
	return b.WriteUint16(uint16(i))
}

func (b *Buffer) ReadUint16() (i uint16, err error) {
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

func (b *Buffer) WriteUint16(i uint16) error {
	m, ok := b.tryGrowByReslice(2)
	if !ok {
		m = b.grow(2)
	}
	binary.BigEndian.PutUint16(b.buf[m:m+2], i)
	return nil
}

func (b *Buffer) ReadInt32() (int32, error) {
	n, err := b.ReadUint32()
	return int32(n), err
}

func (b *Buffer) WriteInt32(i int32) error {
	return b.WriteUint32(uint32(i))
}

func (b *Buffer) ReadUint32() (i uint32, err error) {
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

func (b *Buffer) WriteUint32(i uint32) error {
	m, ok := b.tryGrowByReslice(4)
	if !ok {
		m = b.grow(4)
	}
	binary.BigEndian.PutUint32(b.buf[m:m+4], i)
	return nil
}

func (b *Buffer) ReadInt64() (int64, error) {
	n, err := b.ReadUint64()
	return int64(n), err
}

func (b *Buffer) WriteInt64(i int64) error {
	return b.WriteUint64(uint64(i))
}

func (b *Buffer) ReadUint64() (i uint64, err error) {
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

func (b *Buffer) WriteUint64(i uint64) error {
	m, ok := b.tryGrowByReslice(8)
	if !ok {
		m = b.grow(8)
	}
	binary.BigEndian.PutUint64(b.buf[m:m+8], i)
	return nil
}

func (b *Buffer) ReadBool() (bool, error) {
	c, err := b.ReadByte()
	if err != nil {
		return false, err
	}
	return c == 1, nil
}

func (b *Buffer) WriteBool(v bool) error {
	if v {
		return b.WriteByte(1)
	} else {
		return b.WriteByte(0)
	}
}

func (b *Buffer) ReadVarint16() (int16, error) {
	i, n := binary.Varint(b.buf[b.off:])
	if n == 0 {
		return 0, io.EOF
	}
	if n < 0 || n > MaxVarintLen16 {
		return 0, ErrVarintOverflow
	}
	b.off += n
	return int16(i), nil
}

func (b *Buffer) WriteVarint16(i int16) error {
	var buf [MaxVarintLen16]byte
	n := binary.PutVarint(buf[:], int64(i))

	m, ok := b.tryGrowByReslice(n)
	if !ok {
		m = b.grow(n)
	}

	n = copy(b.buf[m:m+n], buf[:n])
	return nil
}

func (b *Buffer) ReadUvarint16() (uint16, error) {
	i, n := binary.Uvarint(b.buf[b.off:])
	if n == 0 {
		return 0, io.EOF
	}
	if n < 0 || n > MaxVarintLen16 {
		return 0, ErrVarintOverflow
	}
	b.off += n
	return uint16(i), nil
}

func (b *Buffer) WriteUvarint16(i uint16) error {
	var buf [MaxVarintLen16]byte
	n := binary.PutUvarint(buf[:], uint64(i))

	m, ok := b.tryGrowByReslice(n)
	if !ok {
		m = b.grow(n)
	}

	n = copy(b.buf[m:m+n], buf[:n])
	return nil
}

func (b *Buffer) ReadVarint32() (int32, error) {
	i, n := binary.Varint(b.buf[b.off:])
	if n == 0 {
		return 0, io.EOF
	}
	if n < 0 || n > MaxVarintLen32 {
		return 0, ErrVarintOverflow
	}
	b.off += n
	return int32(i), nil
}

func (b *Buffer) WriteVarint32(i int32) error {
	var buf [MaxVarintLen32]byte
	n := binary.PutVarint(buf[:], int64(i))

	m, ok := b.tryGrowByReslice(n)
	if !ok {
		m = b.grow(n)
	}

	n = copy(b.buf[m:m+n], buf[:n])
	return nil
}

func (b *Buffer) ReadUvarint32() (uint32, error) {
	i, n := binary.Uvarint(b.buf[b.off:])
	if n == 0 {
		return 0, io.EOF
	}
	if n < 0 || n > MaxVarintLen32 {
		return 0, ErrVarintOverflow
	}
	b.off += n
	return uint32(i), nil
}

func (b *Buffer) WriteUvarint32(i uint32) error {
	var buf [MaxVarintLen32]byte
	n := binary.PutUvarint(buf[:], uint64(i))

	m, ok := b.tryGrowByReslice(n)
	if !ok {
		m = b.grow(n)
	}

	n = copy(b.buf[m:m+n], buf[:n])
	return nil
}

func (b *Buffer) ReadVarint64() (int64, error) {
	i, n := binary.Varint(b.buf[b.off:])
	if n == 0 {
		return 0, io.EOF
	}
	if n < 0 {
		return 0, ErrVarintOverflow
	}
	b.off += n
	return i, nil
}

func (b *Buffer) WriteVarint64(i int64) error {
	var buf [MaxVarintLen64]byte
	n := binary.PutVarint(buf[:], i)

	m, ok := b.tryGrowByReslice(n)
	if !ok {
		m = b.grow(n)
	}

	n = copy(b.buf[m:m+n], buf[:n])
	return nil
}

func (b *Buffer) ReadUvarint64() (uint64, error) {
	i, n := binary.Uvarint(b.buf[b.off:])
	if n == 0 {
		return 0, io.EOF
	}
	if n < 0 {
		return 0, ErrVarintOverflow
	}
	b.off += n
	return i, nil
}

func (b *Buffer) WriteUvarint64(i uint64) error {
	var buf [MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], i)

	m, ok := b.tryGrowByReslice(n)
	if !ok {
		m = b.grow(n)
	}

	n = copy(b.buf[m:m+n], buf[:n])
	return nil
}

func (b *Buffer) Read(p []byte) (int, error) {
	if len(p) <= 0 {
		return 0, io.ErrShortBuffer
	}

	if b.Readable() == 0 {
		return 0, io.EOF
	}

	n := copy(p, b.buf[b.off:])
	b.off += n
	return n, nil
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	m, ok := b.tryGrowByReslice(len(p))
	if !ok {
		m = b.grow(len(p))
	}
	return copy(b.buf[m:], p), nil
}

func (b *Buffer) ReadString() (string, error) {
	i, n := binary.Varint(b.buf[b.off:])
	if n == 0 {
		return "", io.EOF
	}
	if n < 0 || n > MaxStringLenLen {
		return "", pkg_errors.WithMessage(ErrVarintOverflow, "read length")
	}

	l := int(i)
	if l < 0 {
		return "", buffer.ErrStringLenExceedLimit
	}

	if l+n > b.Readable() {
		return "", io.ErrUnexpectedEOF
	}

	b.off += n
	if l == 0 {
		return "", nil
	}
	s := string(b.buf[b.off : b.off+l])
	b.off += l
	return s, nil
}

func (b *Buffer) WriteString(s string) error {
	l := len(s)
	if l > math.MaxInt32 {
		return buffer.ErrStringLenExceedLimit
	}

	var buf [MaxStringLenLen]byte
	ll := binary.PutVarint(buf[:], int64(l))

	m, ok := b.tryGrowByReslice(ll + l)
	if !ok {
		m = b.grow(ll + l)
	}

	copy(b.buf[m:m+ll], buf[:ll])
	m = m + ll
	copy(b.buf[m:m+l], s)
	return nil
}

func (b *Buffer) ReadFrom(r io.Reader) (int64, error) {
	l := len(b.buf)
	c := cap(b.buf)
	if l >= c {
		return 0, buffer.ErrBufferFull
	}

	return b.ReadFromN(r, c-l)
}

func (b *Buffer) ReadFromN(r io.Reader, n int) (int64, error) {
	if n < 0 {
		panic("bytes.Buffer.ReadFromN: n < 0")
	}

	m := b.grow(n)
	b.buf = b.buf[:m]

	nn, err := r.Read(b.buf[m : m+n])
	if nn < 0 {
		panic("bytes.Buffer.ReadFromN: reader returned negative count from Read")
	}

	b.buf = b.buf[:m+nn]
	return int64(nn), err
}

func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	l := b.Readable()
	if l == 0 {
		return 0, nil
	}

	n, err := w.Write(b.buf[b.off:])
	if n > l {
		panic("bytes.Buffer: invalid Write count")
	}

	b.off += n
	return int64(n), err
}

func (b *Buffer) Peek(n int) ([]byte, error) {
	if n > len(b.buf) {
		return nil, buffer.ErrExceedBufferLimit
	}

	if n > b.Readable() {
		return b.buf[b.off:], io.ErrUnexpectedEOF
	}

	return b.buf[b.off : b.off+n], nil
}

func (b *Buffer) Skip(n int) (skipped int, err error) {
	if n > len(b.buf) {
		return 0, buffer.ErrExceedBufferLimit
	}

	if l := b.Readable(); n > l {
		skipped = l
		err = io.ErrUnexpectedEOF
	} else {
		skipped = n
	}

	b.off += skipped
	return
}

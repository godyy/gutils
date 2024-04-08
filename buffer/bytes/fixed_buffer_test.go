package bytes

import (
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/rs/xid"
)

func TestFixedBuffer(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	b := NewFixedBuffer(128)

	c := byte(rand.Intn(128))
	if err := b.WriteByte(c); err != nil {
		t.Fatalf("write byte %d: %v", c, err)
	} else {
		t.Logf("write byte %d", c)
	}

	u16 := uint16(rand.Intn(math.MaxUint16))
	if err := b.WriteUint16(u16); err != nil {
		t.Fatalf("write uint16 %d: %v", u16, err)
	} else {
		t.Logf("write uint16 %d", u16)
	}

	u32 := uint32(rand.Intn(math.MaxUint32))
	if err := b.WriteUint32(u32); err != nil {
		t.Fatalf("write uint32 %d: %v", u32, err)
	} else {
		t.Logf("write uint32 %d", u32)
	}

	u64 := uint64(rand.Uint64())
	if err := b.WriteUint64(u64); err != nil {
		t.Fatalf("write uint64 %d: %v", u64, err)
	} else {
		t.Logf("write uint64 %d", u64)
	}

	s := xid.New().String()
	if err := b.WriteString(s); err != nil {
		t.Fatalf("write string %s: %v", s, err)
	} else {
		t.Logf("write string %s", s)
	}

	s2 := strings.Repeat("s2", 10)
	if err := b.WriteString(s2); err != nil {
		t.Fatalf("write string2 %s: %v", s2, err)
	} else {
		t.Logf("write string2 %s", s2)
	}

	vMinI16 := int16(math.MinInt16)
	if _, err := b.WriteVarint16(vMinI16); err != nil {
		t.Fatalf("write min varint16 %d: %v", vMinI16, err)
	} else {
		t.Logf("write min varint16 %d", vMinI16)
	}

	vMaxI16 := int16(math.MaxInt16)
	if _, err := b.WriteVarint16(vMaxI16); err != nil {
		t.Fatalf("write max varint16 %d: %v", vMaxI16, err)
	} else {
		t.Logf("write max varint16 %d", vMaxI16)
	}

	vMaxUI16 := uint16(math.MaxUint16)
	if _, err := b.WriteUvarint16(vMaxUI16); err != nil {
		t.Fatalf("write max varuint16 %d: %v", vMaxUI16, err)
	} else {
		t.Logf("write max varuint16 %d", vMaxUI16)
	}

	vMinI32 := int32(math.MinInt32)
	if _, err := b.WriteVarint32(vMinI32); err != nil {
		t.Fatalf("write min varint32 %d: %v", vMinI32, err)
	} else {
		t.Logf("write min varint32 %d", vMinI32)
	}

	vMaxI32 := int32(math.MaxInt32)
	if _, err := b.WriteVarint32(vMaxI32); err != nil {
		t.Fatalf("write max varint32 %d: %v", vMaxI32, err)
	} else {
		t.Logf("write max varint32 %d", vMaxI32)
	}

	vMaxUI32 := uint32(math.MaxUint32)
	if _, err := b.WriteUvarint32(vMaxUI32); err != nil {
		t.Fatalf("write max varuint32 %d: %v", vMaxUI32, err)
	} else {
		t.Logf("write max varuint32 %d", vMaxUI32)
	}

	vMinI64 := int64(math.MinInt64)
	if _, err := b.WriteVarint64(vMinI64); err != nil {
		t.Fatalf("write min varint64 %d: %v", vMinI64, err)
	} else {
		t.Logf("write min varint64 %d", vMinI64)
	}

	vMaxI64 := int64(math.MaxInt64)
	if _, err := b.WriteVarint64(vMaxI64); err != nil {
		t.Fatalf("write max varint64 %d: %v", vMaxI64, err)
	} else {
		t.Logf("write max varint64 %d", vMaxI64)
	}

	vMaxUI64 := uint64(math.MaxUint64)
	if _, err := b.WriteUvarint64(vMaxUI64); err != nil {
		t.Fatalf("write max varuint64 %d: %v", vMaxUI64, err)
	} else {
		t.Logf("write max varuint64 %d", vMaxUI64)
	}

	t.Logf("buffered: %d, available: %d", b.Readable(), b.Writable())

	if cc, err := b.ReadByte(); err != nil {
		t.Fatalf("read byte: %v", err)
	} else if cc != c {
		t.Fatalf("read byte %d != %d", cc, c)
	} else {
		t.Logf("read byte %d == %d", cc, c)
	}

	if u, err := b.ReadUint16(); err != nil {
		t.Fatalf("read uint16: %v", err)
	} else if u != u16 {
		t.Fatalf("read uint16 %d != %d", u, u16)
	} else {
		t.Logf("read uint16 %d == %d", u, u16)
	}

	if u, err := b.ReadUint32(); err != nil {
		t.Fatalf("read uint32: %v", err)
	} else if u != u32 {
		t.Fatalf("read uint32 %d != %d", u, u32)
	} else {
		t.Logf("read uint32 %d == %d", u, u32)
	}

	if u, err := b.ReadUint64(); err != nil {
		t.Fatalf("read uint64: %v", err)
	} else if u != u64 {
		t.Fatalf("read uint64 %d != %d", u, u64)
	} else {
		t.Logf("read uint64 %d == %d", u, u64)
	}

	if ss, err := b.ReadString(); err != nil {
		t.Fatalf("read string: %v", err)
	} else if ss != s {
		t.Fatalf("read string %s != %s", ss, s)
	} else {
		t.Logf("read string %s == %s", ss, s)
	}

	if ss, err := b.ReadString(); err != nil {
		t.Fatalf("read string: %v", err)
	} else if ss != s2 {
		t.Fatalf("read string %s != %s", ss, s2)
	} else {
		t.Logf("read string %s == %s", ss, s2)
	}

	if i, err := b.ReadVarint16(); err != nil {
		t.Fatalf("read min varint16: %v", err)
	} else if i != vMinI16 {
		t.Fatalf("read min varint16 %d != %d", i, vMinI16)
	} else {
		t.Logf("read min varint16 %d == %d", i, vMinI16)
	}

	if i, err := b.ReadVarint16(); err != nil {
		t.Fatalf("read max varint16: %v", err)
	} else if i != vMaxI16 {
		t.Fatalf("read max varint16 %d != %d", i, vMaxI16)
	} else {
		t.Logf("read max varint16 %d == %d", i, vMaxI16)
	}

	if i, err := b.ReadUvarint16(); err != nil {
		t.Fatalf("read max varuint16: %v", err)
	} else if i != vMaxUI16 {
		t.Fatalf("read max varuint16 %d != %d", i, vMaxUI16)
	} else {
		t.Logf("read max varuint16 %d == %d", i, vMaxUI16)
	}

	if i, err := b.ReadVarint32(); err != nil {
		t.Fatalf("read min varint32: %v", err)
	} else if i != vMinI32 {
		t.Fatalf("read min varint32 %d != %d", i, vMinI32)
	} else {
		t.Logf("read min varint32 %d == %d", i, vMinI32)
	}

	if i, err := b.ReadVarint32(); err != nil {
		t.Fatalf("read max varint32: %v", err)
	} else if i != vMaxI32 {
		t.Fatalf("read max varint32 %d != %d", i, vMaxI32)
	} else {
		t.Logf("read max varint32 %d == %d", i, vMaxI32)
	}

	if i, err := b.ReadUvarint32(); err != nil {
		t.Fatalf("read max varuint32: %v", err)
	} else if i != vMaxUI32 {
		t.Fatalf("read max varuint32 %d != %d", i, vMaxUI32)
	} else {
		t.Logf("read max varuint32 %d == %d", i, vMaxUI32)
	}

	if i, err := b.ReadVarint64(); err != nil {
		t.Fatalf("read min varint64: %v", err)
	} else if i != vMinI64 {
		t.Fatalf("read min varint64 %d != %d", i, vMinI64)
	} else {
		t.Logf("read min varint64 %d == %d", i, vMinI64)
	}

	if i, err := b.ReadVarint64(); err != nil {
		t.Fatalf("read max varint64: %v", err)
	} else if i != vMaxI64 {
		t.Fatalf("read max varint64 %d != %d", i, vMaxI64)
	} else {
		t.Logf("read max varint64 %d == %d", i, vMaxI64)
	}

	if i, err := b.ReadUvarint64(); err != nil {
		t.Fatalf("read max varuint64: %v", err)
	} else if i != vMaxUI64 {
		t.Fatalf("read max varuint64 %d != %d", i, vMaxUI64)
	} else {
		t.Logf("read max varuint64 %d == %d", i, vMaxUI64)
	}

	t.Logf("buffered: %d, available: %d", b.Readable(), b.Writable())

	if err := b.WriteUint64(128); err != nil {
		t.Fatalf("write 8 byte: %v", err)
	}
	if n, err := b.Write(make([]byte, 120)); err != nil {
		t.Fatalf("write 120 byte: %v", err)
	} else if n < 120 {
		t.Fatalf("write 120 byte short")
	}
	t.Logf("buffered: %d, available: %d", b.Readable(), b.Writable())
	if _, err := b.ReadUint64(); err != nil {
		t.Fatalf("read 8 byte: %v", err)
	}
	t.Logf("buffered: %d, available: %d", b.Readable(), b.Writable())
	if err := b.WriteUint64(128); err != nil {
		t.Fatalf("read 8 byte at last: %v", err)
	}
	t.Logf("buffered: %d, available: %d", b.Readable(), b.Writable())
}

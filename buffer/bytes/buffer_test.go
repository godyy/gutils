package bytes

import (
	"bytes"
	"github.com/rs/xid"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestBuffer(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	b := NewBufferWithBuf()

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

	vi := int64(math.MaxInt64)
	if err := b.WriteVarint(vi); err != nil {
		t.Fatalf("write varint %d: %v", vi, err)
	} else {
		t.Logf("write varint %d", vi)
	}

	t.Logf("size:%d, cap:%d, reaable: %d, writable: %d", b.Size(), b.Cap(), b.Readable(), b.Writable())

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

	if i, err := b.ReadVarint(); err != nil {
		t.Fatalf("read varint: %v", err)
	} else if i != vi {
		t.Fatalf("read varint %d != %d", i, vi)
	} else {
		t.Logf("read varint %d == %d", i, vi)
	}

	t.Logf("size:%d, cap:%d, reaable: %d, writable: %d", b.Size(), b.Cap(), b.Readable(), b.Writable())

	if err := b.WriteUint64(128); err != nil {
		t.Fatalf("write 8 byte: %v", err)
	}
	_, _ = b.ReadInt32()
	if n, err := b.Write(make([]byte, 124)); err != nil {
		t.Fatalf("write 120 byte: %v", err)
	} else if n < 120 {
		t.Fatalf("write 120 byte short")
	}
	t.Logf("size:%d, cap:%d, reaable: %d, writable: %d", b.Size(), b.Cap(), b.Readable(), b.Writable())
	if _, err := b.ReadUint64(); err != nil {
		t.Fatalf("read 8 byte: %v", err)
	}
	t.Logf("size:%d, cap:%d, reaable: %d, writable: %d", b.Size(), b.Cap(), b.Readable(), b.Writable())
	if err := b.WriteUint64(128); err != nil {
		t.Fatalf("read 8 byte at last: %v", err)
	}
	t.Logf("size:%d, cap:%d, reaable: %d, writable: %d", b.Size(), b.Cap(), b.Readable(), b.Writable())

	b.Reset()
	bs, bc := b.Size(), b.Cap()
	bb := bytes.NewBuffer(nil)
	bb.Write(make([]byte, b.Cap()+50))
	_, _ = b.ReadFrom(bb)
	if b.Size() != bc || bc != b.Cap() {
		t.Fatalf("size=%d cap=%d, after ReadFrom, size=%d cap=%d", bs, bc, b.Size(), b.Cap())
	}

}

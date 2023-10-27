package bytes

import (
	"encoding/binary"
	"errors"
)

// ErrVarintOverflow  varint值溢出
var ErrVarintOverflow = errors.New("bytes: varint overflow")

// MaxVarintLenN is the maximum length of a varint-encoded N-bit integer.
const (
	MaxVarintLen16 = binary.MaxVarintLen16
	MaxVarintLen32 = binary.MaxVarintLen32
	MaxVarintLen64 = binary.MaxVarintLen64
)

const MaxStringLenLen = MaxVarintLen32

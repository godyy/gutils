package bytes

import "errors"

// ErrVarintOverflow  varint值溢出
var ErrVarintOverflow = errors.New("bytes: varint overflow")

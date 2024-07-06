package base

// source: "https://github.com/golang/go/issues/54582"
import (
	"strconv"
	"sync/atomic"
)

type JsonAtomicInt struct {
	atomic.Int64
}

// MarshalText implements the encoding.TextMarshaler interface.
func (x *JsonAtomicInt) MarshalText() ([]byte, error) {
	if x == nil {
		return []byte("<nil>"), nil
	}
	return []byte(strconv.FormatInt(x.Load(), 10)), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (x *JsonAtomicInt) UnmarshalText(text []byte) error {
	n, err := strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		return err
	}
	x.Store(n)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (x *JsonAtomicInt) MarshalJSON() ([]byte, error) {
	if x == nil {
		return []byte("null"), nil
	}
	return x.MarshalText()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (x *JsonAtomicInt) UnmarshalJSON(text []byte) error {
	// Ignore null, like in the main JSON package.
	if string(text) == "null" {
		return nil
	}
	return x.UnmarshalText(text)
}

package log

import "fmt"

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

var (
	ErrOffsetNotFound = fmt.Errorf("offset not found")
	ErrLogIsEmpty     = fmt.Errorf("log is empty")
)

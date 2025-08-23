package log

import (
	"sync"
)

type Config struct {
	InitialCapacity uint64
}

type Log struct {
	mu      sync.RWMutex
	records []Record
}

func NewLog(config Config) *Log {
	var records []Record
	if config.InitialCapacity > 0 {
		records = make([]Record, 0, config.InitialCapacity)
	}
	return &Log{
		records: records,
	}
}

func (l *Log) Append(record Record) (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	record.Offset = uint64(len(l.records))
	l.records = append(l.records, record)

	return record.Offset, nil
}

func (l *Log) Read(offset uint64) (Record, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if offset >= uint64(len(l.records)) {
		return Record{}, ErrOffsetNotFound
	}

	return l.records[offset], nil
}

func (l *Log) HighestOffset() (uint64, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if len(l.records) == 0 {
		return 0, ErrLogIsEmpty
	}

	return uint64(len(l.records) - 1), nil
}

func (l *Log) Truncate(lowest uint64) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if lowest >= uint64(len(l.records)) {
		l.records = []Record{}
		return nil
	}

	l.records = l.records[lowest:]
	return nil
}

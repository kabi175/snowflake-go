package generator

import (
	"errors"
	"time"
)

type Generator struct {
	prev      int64
	offset    int64
	machineID int64
	sequence  int64
}

func NewGenerator(machineID int64, offset int64) (*Generator, error) {
	if machineID < 0 || machineID >= 1024 {
		return nil, errors.New("incompatable machine id")
	}

	if offset < 0 || offset > time.Now().UnixMilli() {
		return nil, errors.New("offset can't be future time")
	}

	return &Generator{
		prev:      offset,
		offset:    offset,
		machineID: machineID << 12,
		sequence:  0,
	}, nil
}

func (g *Generator) Next() (int64, error) {
	curr := time.Now().UnixMilli()

	curr -= g.offset

	next := (curr << 22) | g.machineID<<12 | g.sequence<<0

	if next < 0 {
		return 0, errors.New("negative id genrated")
	}

	if next < g.prev {
		return 0, errors.New("backward time flow")
	}

	g.prev = next
	g.sequence = (g.sequence + 1) % 4096

	return next, nil
}

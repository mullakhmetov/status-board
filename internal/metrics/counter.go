package metrics

import "sync/atomic"

type Counter interface {
	Name() string
	Inc()
	Count() int64
}

func NewDummyCounter() Counter {
	return &DummyCounter{}
}

type DummyCounter struct{}

func (d *DummyCounter) Name() string {
	return ""
}

func (d *DummyCounter) Inc() {}

func (d *DummyCounter) Count() int64 { return 0 }

func NewStandardCounter(name string) Counter {
	return &StandardCounter{name, 0}
}

type StandardCounter struct {
	name  string
	count int64
}

func (c *StandardCounter) Name() string {
	return c.name + " checks"
}

func (c *StandardCounter) Inc() {
	atomic.AddInt64(&c.count, 1)
}

func (c *StandardCounter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

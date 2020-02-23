package metrics

import "sync"

// NewRegistry one of Registry
func NewRegistry(dummy bool) *Registry {
	var initCounterFunc func(string) Counter

	if dummy {
		initCounterFunc = func(name string) Counter {
			return NewDummyCounter()
		}
	} else {
		initCounterFunc = NewStandardCounter
	}
	return &Registry{
		InitCounterFunc: initCounterFunc,
		Counters:        make(map[string]Counter),
		dummy:           dummy,
	}
}

// Registry stores Counters set
type Registry struct {
	lock            sync.Mutex
	InitCounterFunc func(string) Counter
	Counters        map[string]Counter

	dummy bool
}

// AddCounter inits new Counter by name and adds it no Registry
func (r *Registry) AddCounter(name string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.Counters[name] = r.InitCounterFunc(name)

	return
}

// Stats returns all counters `name: values` map
func (r *Registry) Stats() map[string]int64 {
	m := make(map[string]int64)

	for _, v := range r.Counters {
		m[v.Name()] = v.Count()
	}

	return m
}

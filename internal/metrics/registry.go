package metrics

import "sync"

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

type Registry struct {
	lock            sync.Mutex
	InitCounterFunc func(string) Counter
	Counters        map[string]Counter

	dummy bool
}

func (r *Registry) AddCounter(name string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.Counters[name] = r.InitCounterFunc(name)

	return
}

func (r *Registry) Stats() map[string]int64 {
	m := make(map[string]int64)

	for _, v := range r.Counters {
		m[v.Name()] = v.Count()
	}

	return m
}

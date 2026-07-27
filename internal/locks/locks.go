package locks

import (
	"sync"
)

// PathLocks manages a set of mutexes keyed by file paths to serialize concurrent operations on the same path.
type PathLocks struct {
	locks sync.Map
}

// Lock locks the mutex for the given path and returns an unlock function.
// Usage: defer locks.Lock(path)()
func (pl *PathLocks) Lock(path string) (unlock func()) {
	val, _ := pl.locks.LoadOrStore(path, &sync.Mutex{})
	mtx := val.(*sync.Mutex)
	mtx.Lock()
	unlock = func() {
		mtx.Unlock()
	}

	return
}

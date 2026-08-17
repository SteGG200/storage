package handler

import (
	"os"
	"sync"
	"time"
)

// progressTracker broadcasts upload progress to multiple SSE subscribers.
type progressTracker struct {
	mu          sync.RWMutex
	subscribers map[chan int64]struct{}
	current     int64
	totalSize   int64
	modTime     *time.Time
	closed      bool
	tempPath    string
	expireTimer *time.Timer
	cleanupOnce sync.Once
	cleanupFunc func()
}

func newProgressTracker(tempPath string, totalSize int64, modTime *time.Time) *progressTracker {
	return &progressTracker{
		subscribers: make(map[chan int64]struct{}),
		tempPath:    tempPath,
		totalSize:   totalSize,
		modTime:     modTime,
	}
}

// Publish broadcasts the latest written bytes to all active subscribers.
func (t *progressTracker) Publish(written int64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return
	}
	t.current = written

	for ch := range t.subscribers {
		select {
		case ch <- written:
		default:
			// Drain old buffered value if subscriber is slow, then send latest
			select {
			case <-ch:
			default:
			}
			select {
			case ch <- written:
			default:
			}
		}
	}
}

// Subscribe creates a new buffered channel for an SSE subscriber and sends current progress.
func (t *progressTracker) Subscribe() (chan int64, func(), int64, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil, func() {}, 0, false
	}

	ch := make(chan int64, 10)
	ch <- t.current
	t.subscribers[ch] = struct{}{}

	unsubscribe := func() {
		t.mu.Lock()
		defer t.mu.Unlock()
		if t.subscribers != nil {
			delete(t.subscribers, ch)
		}
	}

	return ch, unsubscribe, t.totalSize, true
}

// Close closes all subscriber channels and marks the tracker closed.
func (t *progressTracker) Close() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return
	}
	t.closed = true
	for ch := range t.subscribers {
		close(ch)
	}
	t.subscribers = nil
}

// Cleanup removes the tracker from store, closes channels, and deletes the temporary file.
func (t *progressTracker) Cleanup() {
	t.cleanupOnce.Do(func() {
		if t.cleanupFunc != nil {
			t.cleanupFunc()
		}
		t.Close()
		//#nosec G703
		_ = os.Remove(t.tempPath)
	})
}

type uploadProgressStore struct {
	store sync.Map
}

func (s *uploadProgressStore) Register(path, tempPath string, totalSize int64, modTime *time.Time) *progressTracker {
	t := newProgressTracker(tempPath, totalSize, modTime)
	t.cleanupFunc = func() {
		s.store.Delete(path)
	}
	t.expireTimer = time.AfterFunc(time.Hour, func() {
		t.Cleanup()
	})
	s.store.Store(path, t)
	return t
}

func (s *uploadProgressStore) Subscribe(path string) (chan int64, func(), int64, bool) {
	val, ok := s.store.Load(path)
	if !ok {
		return nil, func() {}, 0, false
	}
	tracker, ok := val.(*progressTracker)
	if !ok {
		return nil, func() {}, 0, false
	}
	return tracker.Subscribe()
}

func (s *uploadProgressStore) Tracker(path string) (*progressTracker, bool) {
	val, ok := s.store.Load(path)
	if !ok {
		return nil, false
	}
	tracker, ok := val.(*progressTracker)
	if !ok {
		return nil, false
	}

	tracker.mu.RLock()
	isClosed := tracker.closed
	tracker.mu.RUnlock()
	if isClosed {
		return nil, false
	}

	if tracker.expireTimer != nil {
		_ = tracker.expireTimer.Stop()
	}
	return tracker, ok
}

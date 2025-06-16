package lock

import "sync"

type SyncGroup[T func(values ...any)] struct{ sw *sync.WaitGroup }

// NEW 实例化
func (SyncGroup[T]) NEW() *SyncGroup[T] { return &SyncGroup[T]{sw: &sync.WaitGroup{}} }

func (my *SyncGroup[T]) Go(fn T) {
	my.sw.Add(1)
	go func() {
		defer my.sw.Done()
		fn()
	}()
}

func (my *SyncGroup[T]) Wait() { my.sw.Wait() }

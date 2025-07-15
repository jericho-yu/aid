package coroutinePool

import (
	"sync"

	"github.com/jericho-yu/aid/array"
)

type (
	CoroutinePool struct {
		sw     sync.WaitGroup
		size   uint // 批次最大执行数
		funcs  *array.AnyArray[func() error]
		wrongs *array.AnyArray[error]
		lock   sync.RWMutex
	}
)

var (
	CoroutinePoolApp CoroutinePool
)

func (*CoroutinePool) New(size uint) *CoroutinePool {
	return &CoroutinePool{sw: sync.WaitGroup{}, size: size, funcs: array.Make[func() error](0), wrongs: array.Make[error](0)}
}

func (my *CoroutinePool) Do(fn func() error) {
	my.funcs.Append(fn)
	my.lock.Lock()
	defer my.lock.Unlock()
}

func (my *CoroutinePool) Close() {
	my.lock.Lock()
	defer my.lock.Unlock()

	if my.funcs.Len() == 0 {
		return
	}

	if my.size == 0 {
		my.size = 1
	}

	chunk := my.funcs.Chunk(int(my.size))
	my.sw.Add(len(chunk))

	for _, funcs := range chunk {
		go func() {
			defer my.sw.Done()
			for _, fn := range funcs {
				go func() { my.wrongs.Append(fn()) }()
			}
		}()
	}

	my.sw.Wait()
}

func (my *CoroutinePool) Wrongs() []error {
	my.lock.Lock()
	defer my.lock.Unlock()
	return my.wrongs.ToSlice()
}

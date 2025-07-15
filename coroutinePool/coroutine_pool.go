package coroutinePool

import (
	"sync"

	"github.com/jericho-yu/aid/array"
)

type (
	// CoroutinePool 协程池
	CoroutinePool struct {
		sw     sync.WaitGroup
		size   uint // 批次最大执行数
		funcs  *array.AnyArray[CoroutinePoolHandle]
		wrongs *array.AnyArray[error]
		lock   sync.RWMutex
	}

	CoroutinePoolHandle func() error
)

var (
	CoroutinePoolApp CoroutinePool
)

// New 实例化：协程池
func (*CoroutinePool) New(size uint) *CoroutinePool {
	return &CoroutinePool{sw: sync.WaitGroup{}, size: size, funcs: array.Make[CoroutinePoolHandle](0), wrongs: array.Make[error](0)}
}

// Do 注册执行方法
func (my *CoroutinePool) Do(fn CoroutinePoolHandle) {
	my.funcs.Append(fn)
	my.lock.Lock()
	defer my.lock.Unlock()
}

// Clean 清空
func (my *CoroutinePool) Clean() *CoroutinePool {
	my.lock.Lock()
	defer my.lock.Unlock()

	my.funcs = array.Make[CoroutinePoolHandle](0)
	my.wrongs = array.Make[error](0)

	return my
}

// Close 关闭并执行
func (my *CoroutinePool) Close() []error {
	my.lock.Lock()
	defer my.lock.Unlock()

	if my.funcs.Len() == 0 {
		return my.wrongs.ToSlice()
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

	return my.wrongs.ToSlice()
}

// Wrongs 获取错误列表
func (my *CoroutinePool) Wrongs() []error {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.wrongs.ToSlice()
}

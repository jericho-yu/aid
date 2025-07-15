package coroutinePool

import "testing"

func Test1(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		names := [...]string{"张三", "李四", "王五", "赵六"}

		pool := CoroutinePoolApp.New(3)
		defer pool.Close()

		for _, name := range names {
			pool.Do(func() error {
				println("协程打印：", name)
				return nil
			})
		}
	})
}

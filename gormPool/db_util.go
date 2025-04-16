package gormPool

import (
	"fmt"

	"gorm.io/gorm"
)

// Finder 查询帮助器
type Finder struct {
	DB    *gorm.DB
	Total int64
}

var FinderApp Finder

// New 实例化：查询帮助器
func (*Finder) New(db *gorm.DB) *Finder { return &Finder{DB: db} }

// Preloads 为查询添加预加载条件。接收一个或多个预加载字段名称作为参数，
// 每个字段名称将被添加到 GORM 的预加载队列中。
// 返回 Finder 实例以支持链式调用。
func (my *Finder) Preloads(preloads ...string) *Finder {
	for _, preload := range preloads {
		my.DB.Preload(preload)
	}

	return my
}

// Find 查询数据
func (my *Finder) Find(ret any) *Finder {
	my.DB.Find(ret)

	return my
}

// Ex 额外操作
func (my *Finder) Ex(fn func(db *gorm.DB)) *Finder {
	fn(my.DB)
	return my
}

// Pagination 分页处理
func (my *Finder) Pagination(page, size int) *Finder {
	if page > 0 && size > 0 {
		if my.DB.Count(&my.Total).Error != nil {
			return my
		}
		my.DB.Limit(size).Offset((page - 1) * size)
	}

	return my
}

// When 当条件满足时执行：where
func (my *Finder) When(condition bool, query any, args ...any) *Finder {
	if condition {
		my.DB.Where(query, args...)
	}

	return my
}

// WhenIn 当条件满足时执行：where in
func (my *Finder) WhenIn(condition bool, query any, args ...any) *Finder {
	if condition {
		my.DB.Where(fmt.Sprintf("%v in (?)", query), args...)
	}

	return my
}

// WhenNotIn 当条件满足时执行：where not in
func (my *Finder) WhenNotIn(condition bool, query any, args ...any) *Finder {
	if condition {
		my.DB.Where(fmt.Sprintf("%v not in (?)", query), args...)
	}

	return my
}

// WhenBetween 当条件满足时执行：where between
func (my *Finder) WhenBetween(condition bool, query any, args ...any) *Finder {
	if condition {
		my.DB.Where(fmt.Sprintf("%v between ? and ?", query), args...)
	}

	return my
}

// WhenNotBetween 当条件满足时执行：where not between
func (my *Finder) WhenNotBetween(condition bool, query any, args ...any) *Finder {
	if condition {
		my.DB.Where(fmt.Sprintf("%v not between ? and ?", query), args...)
	}

	return my
}

// WhenLike 当条件满足时执行：like %?%
func (my *Finder) WhenLike(condition bool, query, arg any) *Finder {
	if condition {
		my.DB.Where(fmt.Sprintf("%v like ?", query), fmt.Sprintf("%%%s%%", arg))
	}

	return my
}

// WhenLikeLeft 当条件满足时执行：like %?
func (my *Finder) WhenLikeLeft(condition bool, query, arg any) *Finder {
	if condition {
		my.DB.Where(fmt.Sprintf("%v like ?", query), fmt.Sprintf("%%%s", arg))
	}

	return my
}

// WhenLikeRight 当条件满足时执行：like ?%
func (my *Finder) WhenLikeRight(condition bool, query, arg any) *Finder {
	if condition {
		my.DB.Where(fmt.Sprintf("%v like ?", query), fmt.Sprintf("%s%%", arg))
	}

	return my
}

// WhenFunc 当条件满足时执行：通过回调执行
func (my *Finder) WhenFunc(condition bool, fn func(db *gorm.DB)) *Finder {
	if condition {
		fn(my.DB)
	}

	return my
}

// Transaction 执行一组数据库事务操作
// 参数 funcs 为需要在事务中执行的函数切片,每个函数接收一个 *gorm.DB 参数
// 如果任一函数执行出错,将回滚整个事务并返回错误
// 所有函数执行成功后提交事务
// 返回 error,nil 表示事务执行成功,非 nil 表示事务执行失败
func (my *Finder) Transaction(funcs ...func(db *gorm.DB)) error {
	my.DB.Begin()

	for _, fn := range funcs {
		fn(my.DB)
		if my.DB.Error != nil {
			my.DB.Rollback()
			return my.DB.Error
		}
	}

	my.DB.Commit()

	return nil
}

package gormPool

import "gorm.io/gorm"

// Finder 查询帮助器
type Finder struct {
	DB    *gorm.DB
	Total int64
}

var FinderApp Finder

// New 实例化：查询帮助器
func (*Finder) New(db *gorm.DB) *Finder { return &Finder{DB: db} }

// FindAndPagination 查询多条数据并自动分页
func (my *Finder) FindAndPagination(page, size int, ret any) *Finder {
	if my.Pagination(page, size).DB.Error != nil {
		return my
	}

	my.DB.Find(ret)

	return my
}

// Pagination 分页处理
func (my *Finder) Pagination(page, size int) *Finder {
	if page > 0 && size > 0 {
		if my.DB.Count(&my.Total).Error != nil {
			return my
		}
		my.DB = my.DB.Limit(size).Offset((page - 1) * size)
	}

	return my
}

// When 条件判断
func (my *Finder) When(condition bool, query any, args ...any) *Finder {
	if condition {
		my.DB = my.DB.Where(query, args...)
	}

	return my
}

// WhenFunc 条件判断：通过回调执行
func (my *Finder) WhenFunc(condition bool, fn func(db *gorm.DB) *gorm.DB) *Finder {
	if condition {
		my.DB = fn(my.DB)
	}

	return my
}

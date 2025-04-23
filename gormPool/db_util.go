package gormPool

import (
	"fmt"

	"gorm.io/gorm"
)

type (
	// Finder 查询帮助器
	Finder struct {
		DB    *gorm.DB
		Total int64
	}
)

var (
	FinderApp Finder
)

// New 实例化：查询帮助器
func (*Finder) New(db *gorm.DB) *Finder { return &Finder{DB: db} }

// Find 查询数据
func (my *Finder) Find(ret any) *Finder {
	my.DB.Find(ret)

	return my
}

// Ex 额外操作
func (my *Finder) Ex(funcs ...func(db *gorm.DB)) *Finder {
	if len(funcs) > 0 {
		for _, fn := range funcs {
			fn(my.DB)
		}
	}

	return my
}

// TryPagination 尝试分页
func (my *Finder) TryPagination(page, size int) *Finder {
	if page > 0 && size > 0 {
		if my.DB.Count(&my.Total).Error != nil {
			return my
		}
		my.DB.Limit(size).Offset((page - 1) * size)
	}

	return my
}

// TryOrder 尝试排序
func (my *Finder) TryOrder(orders ...string) *Finder {
	for _, order := range orders {
		my.DB.Order(order)
	}

	return my
}

// TryPreloads 尝试深度查询
func (my *Finder) TryPreload(preloads ...string) *Finder {
	for _, preload := range preloads {
		my.DB.Preload(preload)
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

// TryQueryFromMap 从map中解析参数并查询
func (my *Finder) TryQueryFromMap(values map[string][]any) *Finder {
	for key, value := range values {
		var (
			ok       = false
			operator string
		)

		if operator, ok = value[0].(string); !ok {
			continue
		}

		switch operator {
		case "alias":
			tableAlias := fmt.Sprintf("%s as %s", key, value[1])
			my.DB.Table(tableAlias)
		case "=", ">", "<", "!=", "<=", ">=", "<>":
			my.DB.Where(fmt.Sprintf("%s %s ?", key, operator), value[1])
		case "in", "not in":
			my.DB.Where(fmt.Sprintf("%s %s (?)", key, operator), value[1])
		case "between", "not between":
			my.DB.Where(fmt.Sprintf("%s %s ? and ?", key, operator), value[1:]...)
		case "like":
			my.DB.Where(fmt.Sprintf("%s like ?", key), fmt.Sprintf("%%%s%%", value[1]))
		case "like%":
			my.DB.Where(fmt.Sprintf("%s like ?", key), fmt.Sprintf("%s%%", value[1]))
		case "%like":
			my.DB.Where(fmt.Sprintf("%s like ?", key), fmt.Sprintf("%%%s", value[1]))
		case "join":
			my.DB.Joins(key, value[1:]...)
		case "raw":
			my.DB.Where(key, value[1:]...)
		case "distant":
			distantList := make([]any, len(value[1:]))
			for i, v := range value[1:] {
				distantList[i] = fmt.Sprintf("%s.%v", key, v)
			}
			my.DB.Distinct(distantList...)
		}
	}

	return my
}

// TryAutoQuery 自动填充查询条件和预加载字段
func (my *Finder) TryAutoFind(queries map[string][]any, preloads []string, orders []string, page, size int, ret any) *Finder {
	my.TryQueryFromMap(queries).TryPreload(preloads...).TryPagination(page, size).TryOrder(orders...).Find(ret)

	return my
}

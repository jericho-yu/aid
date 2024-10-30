package excel

import (
	"errors"
	"fmt"

	"github.com/jericho-yu/aid/array"

	"github.com/xuri/excelize/v2"
)

// Row Excel行
type Row struct {
	Err       error
	cells     *array.AnyArray[*Cell]
	rowNumber uint64
}

// NewRow 构造函数
func NewRow() *Row {
	return &Row{}
}

// GetCells 获取单元格组
func (r *Row) GetCells() *array.AnyArray[*Cell] {
	return r.cells
}

// SetCells 设置单元格组
func (r *Row) SetCells(cells []*Cell) *Row {
	if r.GetRowNumber() == 0 {
		r.Err = errors.New("行标必须大于0")
		return r
	}

	for colNumber, cell := range cells {
		if colText, err := excelize.ColumnNumberToName(colNumber + 1); err != nil {
			panic(fmt.Errorf("列索引转列文字失败：%d，%d", r.GetRowNumber(), colNumber+1))
		} else {
			cell.SetCoordinate(fmt.Sprintf("%s%d", colText, r.GetRowNumber()))
		}
	}
	r.cells = array.NewAnyArray[*Cell](cells)

	return r
}

// GetRowNumber 获取行标
func (r *Row) GetRowNumber() uint64 {
	return r.rowNumber
}

// SetRowNumber 设置行标
func (r *Row) SetRowNumber(rowNumber uint64) *Row {
	r.rowNumber = rowNumber
	return r
}

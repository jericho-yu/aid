package excel

import (
	"errors"
	"fmt"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/jericho-yu/aid/array"
	"github.com/jericho-yu/aid/dict"
	"github.com/xuri/excelize/v2"
)

// Reader Excel读取器
type Reader struct {
	Err         error
	data        *dict.AnyOrderlyDict[uint64, *array.AnyArray[string]]
	excel       *excelize.File
	sheetName   string
	originalRow int
	finishedRow int
	titleRow    int
	titles      *array.AnyArray[string]
}

// NewReader 构造函数
func NewReader() *Reader {
	return &Reader{data: dict.MakeAnyOrderlyDict[uint64, *array.AnyArray[string]](0)}
}

// AutoRead 自动读取（默认第一行是表头，从第二行开始，默认Sheet名称为：Sheet1）
func (r *Reader) AutoRead(filename ...any) *Reader {
	return r.
		OpenFile(filename...).
		SetOriginalRow(2).
		SetTitleRow(1).
		SetSheetName("Sheet1").
		ReadTitle().
		Read()
}

// AutoReadBySheetName 自动读取（默认第一行是表头，从第二行开始）
func (r *Reader) AutoReadBySheetName(sheetName string, filename ...any) *Reader {
	return r.
		OpenFile(filename...).
		SetOriginalRow(2).
		SetTitleRow(1).
		SetSheetName(sheetName).
		ReadTitle().
		Read()
}

// Data 获取数据：有序字典
func (r *Reader) Data() *dict.AnyOrderlyDict[uint64, *array.AnyArray[string]] {
	return r.data
}

// DataWithTitle 获取数据：带有title的有序字典
func (r *Reader) DataWithTitle() (*dict.AnyOrderlyDict[uint64, *dict.AnyDict[string, string]], error) {
	newDict := dict.MakeAnyOrderlyDict[uint64, *dict.AnyDict[string, string]](0)

	for idx, value := range r.data.All() {
		newMap, err := dict.Zip[string, string](r.titles.All(), value.Value.All())
		if err != nil {
			return nil, err
		}
		newDict.SetByKey(uint64(idx), dict.NewAnyDict[string, string](newMap))
	}

	return newDict, nil
}

// SetDataByRow 设置单行数据
func (r *Reader) SetDataByRow(rowNumber uint64, data []string) *Reader {
	r.data.SetByKey(rowNumber, array.NewAnyArray(data))
	return r
}

// GetSheetName 获取工作表名称
func (r *Reader) GetSheetName() string {
	return r.sheetName
}

// SetSheetName 设置工作表名称
func (r *Reader) SetSheetName(sheetName string) *Reader {
	r.sheetName = sheetName
	return r
}

// GetOriginalRow 获取读取起始行
func (r *Reader) GetOriginalRow() int {
	return r.originalRow
}

// SetOriginalRow 设置读取起始行
func (r *Reader) SetOriginalRow(originalRow int) *Reader {
	r.originalRow = originalRow - 1
	return r
}

// GetFinishedRow 获取读取终止行
func (r *Reader) GetFinishedRow() int {
	return r.finishedRow
}

// SetFinishedRow 设置读取终止行
func (r *Reader) SetFinishedRow(finishedRow int) *Reader {
	r.finishedRow = finishedRow - 1
	return r
}

// GetTitleRow 获取表头行
func (r *Reader) GetTitleRow() int {
	return r.titleRow
}

// SetTitleRow 设置表头行
func (r *Reader) SetTitleRow(titleRow int) *Reader {
	r.titleRow = titleRow - 1
	return r
}

// GetTitle 获取表头
func (r *Reader) GetTitle() *array.AnyArray[string] {
	return r.titles
}

// SetTitle 设置表头
func (r *Reader) SetTitle(titles []string) *Reader {
	if len(titles) == 0 {
		r.Err = errors.New("表头不能为空")
		return r
	}
	r.titles = array.NewAnyArray[string](titles)
	return r
}

// OpenFile 打开文件
func (r *Reader) OpenFile(filename ...any) *Reader {
	if filename[0].(string) == "" {
		r.Err = errors.New("文件名不能为空")
		return r
	}
	f, err := excelize.OpenFile(fmt.Sprintf(filename[0].(string), filename[1:]...))
	if err != nil {
		r.Err = fmt.Errorf("打开文件错误：%s", err.Error())
		return r
	}
	r.excel = f

	defer func(r *Reader) {
		if err = r.excel.Close(); err != nil {
			r.Err = errors.New("文件关闭错误")
		}
	}(r)

	r.SetTitleRow(1)
	r.SetOriginalRow(2)
	r.data = dict.MakeAnyOrderlyDict[uint64, *array.AnyArray[string]](0)

	return r
}

// ReadTitle 读取表头
func (r *Reader) ReadTitle() *Reader {
	if r.GetSheetName() == "" {
		r.Err = errors.New("未设置工作表名称")
		return r
	}

	if rows, err := r.excel.GetRows(r.GetSheetName()); err != nil {
		panic(fmt.Errorf("读取表头错误：%s", err.Error()))
	} else {
		r.SetTitle(rows[r.GetTitleRow()])
	}

	return r
}

// Read 读取Excel
func (r *Reader) Read() *Reader {
	if r.GetSheetName() == "" {
		r.Err = errors.New("未设置工作表名称")
		return r
	}

	if rows, err := r.excel.GetRows(r.GetSheetName()); err != nil {
		r.Err = errors.New("读取数据错误：%s")
		return r
	} else {
		if r.finishedRow == 0 {
			for rowNumber, values := range rows[r.GetOriginalRow():] {
				r.SetDataByRow(uint64(rowNumber), values)
			}
		} else {
			for rowNumber, values := range rows[r.GetOriginalRow():r.GetFinishedRow()] {
				r.SetDataByRow(uint64(rowNumber), values)
			}
		}
	}

	return r
}

// ToDataFrameDefaultType 获取DataFrame类型数据 通过Excel表头自定义数据类型
func (r *Reader) ToDataFrameDefaultType() dataframe.DataFrame {
	titleWithType := make(map[string]series.Type)
	for _, title := range r.GetTitle().All() {
		titleWithType[title] = series.String
	}

	return r.ToDataFrame(titleWithType)
}

// ToDataFrame 获取DataFrame类型数据
func (r *Reader) ToDataFrame(titleWithType map[string]series.Type) dataframe.DataFrame {
	if r.GetSheetName() == "" {
		panic(errors.New("未设置工作表名称"))
	}

	var _content [][]string

	if rows, err := r.excel.GetRows(r.GetSheetName()); err != nil {
		panic(errors.New("读取数据错误"))
	} else {
		if r.finishedRow == 0 {
			_content = rows[r.GetTitleRow():]
		} else {
			_content = rows[r.GetTitleRow():r.GetFinishedRow()]
		}
	}

	return dataframe.LoadRecords(
		_content,
		dataframe.DetectTypes(false),
		dataframe.DefaultType(series.String),
		dataframe.WithTypes(titleWithType),
	)
}

// ToDataFrameDetectType 获取DataFrame类型数据 通过自动探寻数据类型
func (r *Reader) ToDataFrameDetectType() dataframe.DataFrame {
	if r.GetSheetName() == "" {
		panic(errors.New("未设置工作表名称"))
	}

	var _content [][]string

	if rows, err := r.excel.GetRows(r.GetSheetName()); err != nil {
		panic(errors.New("读取数据错误"))
	} else {
		if r.finishedRow == 0 {
			_content = rows[r.GetTitleRow():]
		} else {
			_content = rows[r.GetTitleRow():r.GetFinishedRow()]
		}
	}

	return dataframe.LoadRecords(
		_content,
		dataframe.DetectTypes(true),
		dataframe.DefaultType(series.String),
	)
}

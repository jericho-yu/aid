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
func (my *Reader) AutoRead(filename ...any) *Reader {
	return my.
		OpenFile(filename...).
		SetOriginalRow(2).
		SetTitleRow(1).
		SetSheetName("Sheet1").
		ReadTitle().
		Read()
}

// AutoReadBySheetName 自动读取（默认第一行是表头，从第二行开始）
func (my *Reader) AutoReadBySheetName(sheetName string, filename ...any) *Reader {
	return my.
		OpenFile(filename...).
		SetOriginalRow(2).
		SetTitleRow(1).
		SetSheetName(sheetName).
		ReadTitle().
		Read()
}

// Data 获取数据：有序字典
func (my *Reader) Data() *dict.AnyOrderlyDict[uint64, *array.AnyArray[string]] {
	return my.data
}

// DataWithTitle 获取数据：带有title的有序字典
func (my *Reader) DataWithTitle() (*dict.AnyOrderlyDict[uint64, *dict.AnyDict[string, string]], error) {
	newDict := dict.MakeAnyOrderlyDict[uint64, *dict.AnyDict[string, string]](0)

	for idx, value := range my.data.All() {
		newMap, err := dict.Zip[string, string](my.titles.All(), value.Value.All())
		if err != nil {
			return nil, err
		}
		newDict.SetByKey(uint64(idx), dict.NewAnyDict[string, string](newMap))
	}

	return newDict, nil
}

// SetDataByRow 设置单行数据
func (my *Reader) SetDataByRow(rowNumber uint64, data []string) *Reader {
	my.data.SetByKey(rowNumber, array.NewAnyArray(data))
	return my
}

// GetSheetName 获取工作表名称
func (my *Reader) GetSheetName() string {
	return my.sheetName
}

// SetSheetName 设置工作表名称
func (my *Reader) SetSheetName(sheetName string) *Reader {
	my.sheetName = sheetName
	return my
}

// GetOriginalRow 获取读取起始行
func (my *Reader) GetOriginalRow() int {
	return my.originalRow
}

// SetOriginalRow 设置读取起始行
func (my *Reader) SetOriginalRow(originalRow int) *Reader {
	my.originalRow = originalRow - 1
	return my
}

// GetFinishedRow 获取读取终止行
func (my *Reader) GetFinishedRow() int {
	return my.finishedRow
}

// SetFinishedRow 设置读取终止行
func (my *Reader) SetFinishedRow(finishedRow int) *Reader {
	my.finishedRow = finishedRow - 1
	return my
}

// GetTitleRow 获取表头行
func (my *Reader) GetTitleRow() int {
	return my.titleRow
}

// SetTitleRow 设置表头行
func (my *Reader) SetTitleRow(titleRow int) *Reader {
	my.titleRow = titleRow - 1
	return my
}

// GetTitle 获取表头
func (my *Reader) GetTitle() *array.AnyArray[string] {
	return my.titles
}

// SetTitle 设置表头
func (my *Reader) SetTitle(titles []string) *Reader {
	if len(titles) == 0 {
		my.Err = errors.New("表头不能为空")
		return my
	}
	my.titles = array.NewAnyArray[string](titles)
	return my
}

// OpenFile 打开文件
func (my *Reader) OpenFile(filename ...any) *Reader {
	if filename[0].(string) == "" {
		my.Err = errors.New("文件名不能为空")
		return my
	}
	f, err := excelize.OpenFile(fmt.Sprintf(filename[0].(string), filename[1:]...))
	if err != nil {
		my.Err = fmt.Errorf("打开文件错误：%s", err.Error())
		return my
	}
	my.excel = f

	defer func(r *Reader) {
		if err = r.excel.Close(); err != nil {
			r.Err = errors.New("文件关闭错误")
		}
	}(my)

	my.SetTitleRow(1)
	my.SetOriginalRow(2)
	my.data = dict.MakeAnyOrderlyDict[uint64, *array.AnyArray[string]](0)

	return my
}

// ReadTitle 读取表头
func (my *Reader) ReadTitle() *Reader {
	if my.GetSheetName() == "" {
		my.Err = errors.New("未设置工作表名称")
		return my
	}

	if rows, err := my.excel.GetRows(my.GetSheetName()); err != nil {
		panic(fmt.Errorf("读取表头错误：%s", err.Error()))
	} else {
		my.SetTitle(rows[my.GetTitleRow()])
	}

	return my
}

// Read 读取Excel
func (my *Reader) Read() *Reader {
	if my.GetSheetName() == "" {
		my.Err = errors.New("未设置工作表名称")
		return my
	}

	if rows, err := my.excel.GetRows(my.GetSheetName()); err != nil {
		my.Err = errors.New("读取数据错误：%s")
		return my
	} else {
		if my.finishedRow == 0 {
			for rowNumber, values := range rows[my.GetOriginalRow():] {
				my.SetDataByRow(uint64(rowNumber), values)
			}
		} else {
			for rowNumber, values := range rows[my.GetOriginalRow():my.GetFinishedRow()] {
				my.SetDataByRow(uint64(rowNumber), values)
			}
		}
	}

	return my
}

// ToDataFrameDefaultType 获取DataFrame类型数据 通过Excel表头自定义数据类型
func (my *Reader) ToDataFrameDefaultType() dataframe.DataFrame {
	titleWithType := make(map[string]series.Type)
	for _, title := range my.GetTitle().All() {
		titleWithType[title] = series.String
	}

	return my.ToDataFrame(titleWithType)
}

// ToDataFrame 获取DataFrame类型数据
func (my *Reader) ToDataFrame(titleWithType map[string]series.Type) dataframe.DataFrame {
	if my.GetSheetName() == "" {
		panic(errors.New("未设置工作表名称"))
	}

	var _content [][]string

	if rows, err := my.excel.GetRows(my.GetSheetName()); err != nil {
		panic(errors.New("读取数据错误"))
	} else {
		if my.finishedRow == 0 {
			_content = rows[my.GetTitleRow():]
		} else {
			_content = rows[my.GetTitleRow():my.GetFinishedRow()]
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
func (my *Reader) ToDataFrameDetectType() dataframe.DataFrame {
	if my.GetSheetName() == "" {
		panic(errors.New("未设置工作表名称"))
	}

	var _content [][]string

	if rows, err := my.excel.GetRows(my.GetSheetName()); err != nil {
		panic(errors.New("读取数据错误"))
	} else {
		if my.finishedRow == 0 {
			_content = rows[my.GetTitleRow():]
		} else {
			_content = rows[my.GetTitleRow():my.GetFinishedRow()]
		}
	}

	return dataframe.LoadRecords(
		_content,
		dataframe.DetectTypes(true),
		dataframe.DefaultType(series.String),
	)
}

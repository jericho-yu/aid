package excel

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/xuri/excelize/v2"
)

// Writer Excel写入器
type Writer struct {
	Err       error
	filename  string
	excel     *excelize.File
	sheetName string
}

// NewWriter 初始化
func NewWriter(filename ...any) *Writer {
	ins := &Writer{}
	if filename[0].(string) == "" {
		ins.Err = errors.New("文件名不能为空")
		return ins
	}
	ins.filename = fmt.Sprintf(filename[0].(string), filename[1:]...)
	ins.excel = excelize.NewFile()

	return ins
}

// GetFilename 获取文件名
func (r *Writer) GetFilename() string {
	return r.filename
}

// SetFilename 设置文件名
func (r *Writer) SetFilename(filename string) *Writer {
	r.filename = filename
	return r
}

// CreateSheet 创建工作表
func (r *Writer) CreateSheet(sheetName string) *Writer {
	if sheetName == "" {
		r.Err = errors.New("工作表名称不能为空")
		return r
	}
	sheetIndex := r.excel.NewSheet(sheetName)
	r.excel.SetActiveSheet(sheetIndex)
	r.sheetName = r.excel.GetSheetName(sheetIndex)

	return r
}

// ActiveSheetByName 选择工作表（根据名称）
func (r *Writer) ActiveSheetByName(sheetName string) *Writer {
	if sheetName == "" {
		r.Err = errors.New("工作表名称不能为空")
		return r
	}
	sheetIndex := r.excel.GetSheetIndex(sheetName)
	r.excel.SetActiveSheet(sheetIndex)
	r.sheetName = sheetName

	return r
}

// ActiveSheetByIndex 选择工作表（根据编号）
func (r *Writer) ActiveSheetByIndex(sheetIndex int) *Writer {
	if sheetIndex < 0 {
		r.Err = errors.New("工作表索引不能小于0")
		return r
	}
	r.excel.SetActiveSheet(sheetIndex)
	r.sheetName = r.excel.GetSheetName(sheetIndex)
	return r
}

// SetSheetName 设置sheet名称
func (r *Writer) SetSheetName(sheetName string) *Writer {
	r.excel.SetSheetName(r.sheetName, sheetName)
	r.sheetName = sheetName
	return r
}

// setStyleFont 设置字体
func (r *Writer) setStyleFont(cell *Cell) {
	fill := excelize.Fill{Type: "pattern", Pattern: 0, Color: []string{""}}
	if cell.GetPatternRgb() != "" {
		fill.Pattern = 1
		fill.Color[0] = cell.GetPatternRgb()
	}

	var borders = make([]excelize.Border, 0)
	if cell.GetBorder().Len() > 0 {
		for _, border := range cell.GetBorder().All() {
			borders = append(borders, excelize.Border{
				Type:  border.Type,
				Color: border.Rgb,
				Style: border.Style,
			})
		}
	}

	if style, err := r.excel.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   cell.GetFontBold(),
			Italic: cell.GetFontItalic(),
			Family: cell.GetFontFamily(),
			Size:   cell.GetFontSize(),
			Color:  cell.GetFontRgb(),
		},
		Alignment: &excelize.Alignment{
			WrapText: cell.GetWrapText(),
		},
		Fill:   fill,
		Border: borders,
	}); err != nil {
		r.Err = fmt.Errorf("设置字体错误：%s", cell.GetCoordinate())
	} else {
		r.Err = r.excel.SetCellStyle(r.sheetName, cell.GetCoordinate(), cell.GetCoordinate(), style)
	}
}

// SetColumnWidthByIndex 设置单列宽：通过列索引
func (r *Writer) SetColumnWidthByIndex(col int, width float64) *Writer {
	r.SetColumnsWidthByIndex(col, col, width)
	return r
}

// SetColumnWidthByText 设置单列宽：通过列名称
func (r *Writer) SetColumnWidthByText(col string, width float64) *Writer {
	r.SetColumnsWidthByText(col, col, width)
	return r
}

// SetColumnsWidthByText 设置多列宽：通过列索引
func (r *Writer) SetColumnsWidthByIndex(startCol, endCol int, width float64) *Writer {
	startColText, err := ColumnNumberToText(startCol)
	if err != nil {
		r.Err = fmt.Errorf("设置列宽错误：%s", err)
	}
	endColText, err := ColumnNumberToText(endCol)
	if err != nil {
		r.Err = fmt.Errorf("设置列宽错误：%s", err)
	}
	if err = r.excel.SetColWidth(r.sheetName, startColText, endColText, width); err != nil {
		r.Err = fmt.Errorf("设置列宽错误：%s", err)
	}
	return r
}

// SetColumnsWidthByText 设置多列宽：通过列名称
func (r *Writer) SetColumnsWidthByText(startCol, endCol string, width float64) *Writer {
	if err := r.excel.SetColWidth(r.sheetName, startCol, endCol, width); err != nil {
		r.Err = fmt.Errorf("设置列宽错误：%s", err)
	}
	return r
}

// SetRows 设置行数据
func (r *Writer) SetRows(excelRows []*Row) *Writer {
	for _, row := range excelRows {
		r.AddRow(row)
	}
	return r
}

// AddRow 增加一行行数据
func (r *Writer) AddRow(excelRow *Row) *Writer {
	for _, cell := range excelRow.GetCells().All() {
		r.Err = r.excel.SetCellValue(r.sheetName, cell.GetCoordinate(), cell.GetContent())
		switch cell.GetContentType() {
		case CellContentTypeFormula:
			if err := r.excel.SetCellFormula(r.sheetName, cell.GetCoordinate(), cell.GetContent().(string)); err != nil {
				r.Err = fmt.Errorf("写入数据错误（公式）%s %s：%v", cell.GetCoordinate(), cell.GetContent(), err.Error())
				return r
			}
		case CellContentTypeAny:
			if err := r.excel.SetCellValue(r.sheetName, cell.GetCoordinate(), cell.GetContent()); err != nil {
				r.Err = fmt.Errorf("写入ExcelCell（任意） %s %s：%v", cell.GetCoordinate(), cell.GetContent(), err.Error())
				return r
			}
		case CellContentTypeInt:
			if err := r.excel.SetCellInt(r.sheetName, cell.GetCoordinate(), cell.GetContent().(int)); err != nil {
				r.Err = fmt.Errorf("写入ExcelCell（整数） %s %s：%v", cell.GetCoordinate(), cell.GetContent(), err.Error())
				return r
			}
		case CellContentTypeFloat64:
			if err := r.excel.SetCellFloat(r.sheetName, cell.GetCoordinate(), cell.GetContent().(float64), 2, 64); err != nil {
				r.Err = fmt.Errorf("写入ExcelCell（浮点数） %s %s：%v", cell.GetCoordinate(), cell.GetContent(), err.Error())
				return r
			}
		case CellContentTypeBool:
			if err := r.excel.SetCellBool(r.sheetName, cell.GetCoordinate(), cell.GetContent().(bool)); err != nil {
				r.Err = fmt.Errorf("写入ExcelCell（布尔） %s %s：%v", cell.GetCoordinate(), cell.GetContent(), err.Error())
				return r
			}
		case CellContentTypeTime:
			if err := r.excel.SetCellValue(r.sheetName, cell.GetCoordinate(), cell.GetContent().(time.Time)); err != nil {
				r.Err = fmt.Errorf("写入ExcelCell（时间） %s %s：%v", cell.GetCoordinate(), cell.GetContent(), err.Error())
			}
		}
		r.setStyleFont(cell)
	}

	return r
}

// SetTitleRow 设置标题行
func (r *Writer) SetTitleRow(titles []string, rowNumber uint64) *Writer {
	var (
		titleRow   *Row
		titleCells = make([]*Cell, len(titles))
	)

	if len(titles) > 0 {
		for idx, title := range titles {
			titleCells[idx] = NewCellAny(title)
		}

		titleRow = NewRow().SetRowNumber(rowNumber).SetCells(titleCells)

		r.AddRow(titleRow)
	}

	return r
}

// Save 保存文件
func (r *Writer) Save() error {
	if r.filename == "" {
		return errors.New("未设置文件名")
	}
	return r.excel.SaveAs(r.filename)
}

// Download 下载Excel
func (r *Writer) Download(writer http.ResponseWriter) error {
	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape(r.filename)))
	writer.Header().Set("Content-Transfer-Encoding", "binary")
	writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	return r.excel.Write(writer)
}

// GetExcelizeFile 获取excelize文件对象
func (r *Writer) GetExcelizeFile() *excelize.File {
	return r.excel
}

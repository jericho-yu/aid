package excel

import (
	"time"

	"github.com/jericho-yu/aid/array"
)

type (
	// CellContentType 单元格内容类型
	CellContentType string

	// Cell Excel单元格
	Cell struct {
		content                                                                                                              any
		contentType                                                                                                          CellContentType
		coordinate                                                                                                           string
		fontRgb                                                                                                              string
		patternRgb                                                                                                           string
		fontBold                                                                                                             bool
		fontItalic                                                                                                           bool
		fontFamily                                                                                                           string
		fontSize                                                                                                             float64
		borderTopRgb, borderBottomRgb, borderLeftRgb, borderRightRgb, borderDiagonalUpRgb, borderDiagonalDownRgb             string
		borderTopStyle, borderBottomStyle, borderLeftStyle, borderRightStyle, borderDiagonalUpStyle, borderDiagonalDownStyle int
		wrapText                                                                                                             bool
	}

	// border 单元格边框
	border struct {
		Type  string
		Rgb   string
		Style int
	}
)

const (
	CellContentTypeAny     CellContentType = "any"
	CellContentTypeFormula CellContentType = "formula"
	CellContentTypeInt     CellContentType = "int"
	CellContentTypeFloat64 CellContentType = "float64"
	CellContentTypeBool    CellContentType = "bool"
	CellContentTypeTime    CellContentType = "time"
)

// NewCellAny 实例化：任意值
func NewCellAny(content any) *Cell {
	return &Cell{content: content, contentType: CellContentTypeAny}
}

// NewCellInt 实例化：整数
func NewCellInt(content any) *Cell {
	return &Cell{content: content, contentType: CellContentTypeInt}
}

// NewCellFloat64 实例化：浮点
func NewCellFloat64(content any) *Cell {
	return &Cell{content: content, contentType: CellContentTypeFloat64}
}

// NewCellBool 实例化：布尔
func NewCellBool(content any) *Cell {
	return &Cell{content: content, contentType: CellContentTypeBool}
}

// NewCellTime 实例化：时间
func NewCellTime(content time.Time) *Cell {
	return &Cell{content: content, contentType: CellContentTypeTime}
}

// NewCellFormula 实例化：公式
func NewCellFormula(content string) *Cell {
	return &Cell{content: content, contentType: CellContentTypeFormula}
}

// GetBorder 获取边框
func (r *Cell) GetBorder() *array.AnyArray[border] {
	borders := array.MakeAnyArray[border](0)

	if r.borderTopRgb != "" {
		borders.Append(border{
			Type:  "top",
			Rgb:   r.borderTopRgb,
			Style: r.borderTopStyle,
		})
	}

	if r.borderBottomRgb != "" {
		borders.Append(border{
			Type:  "bottom",
			Rgb:   r.borderBottomRgb,
			Style: r.borderBottomStyle,
		})
	}

	if r.borderLeftRgb != "" {
		borders.Append(border{
			Type:  "left",
			Rgb:   r.borderLeftRgb,
			Style: r.borderLeftStyle,
		})
	}

	if r.borderRightRgb != "" {
		borders.Append(border{
			Type:  "right",
			Rgb:   r.borderRightRgb,
			Style: r.borderRightStyle,
		})
	}

	if r.borderDiagonalUpRgb != "" {
		borders.Append(border{
			Type:  "diagonalUp",
			Rgb:   r.borderDiagonalUpRgb,
			Style: r.borderDiagonalUpStyle,
		})
	}

	if r.borderDiagonalDownRgb != "" {
		borders.Append(border{
			Type:  "diagonalDown",
			Rgb:   r.borderDiagonalDownRgb,
			Style: r.borderDiagonalDownStyle,
		})
	}

	return borders
}

// SetWrapText 设置自动换行
func (r *Cell) SetWrapText(wrapText bool) *Cell {
	r.wrapText = wrapText
	return r
}

// GetWrapText 获取自动换行
func (r *Cell) GetWrapText() bool {
	return r.wrapText
}

// SetBorderSurrounding 设置四周边框
func (r *Cell) SetBorderSurrounding(borderRgb string, borderStyle int, condition bool) *Cell {
	if condition {
		r.borderTopRgb = borderRgb
		r.borderBottomRgb = borderRgb
		r.borderLeftRgb = borderRgb
		r.borderRightRgb = borderRgb
		r.borderTopStyle = borderStyle
		r.borderBottomStyle = borderStyle
		r.borderLeftStyle = borderStyle
		r.borderRightStyle = borderStyle
	}

	return r
}

// SetBorderSurroundingFunc 设置四周边框 函数
func (r *Cell) SetBorderSurroundingFunc(condition func() (string, int, bool)) *Cell {
	if condition != nil {
		r.SetBorderSurrounding(condition())
	}
	return r
}

// SetBorderTopRgb 设置边框颜色：上
func (r *Cell) SetBorderTopRgb(borderTopRgb string, condition bool) *Cell {
	if condition {
		r.borderTopRgb = borderTopRgb
	}
	return r
}

// SetBorderTopRbgFunc 设置边框颜色：上 函数
func (r *Cell) SetBorderTopRbgFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		r.SetBorderTopRgb(condition())
	}
	return r
}

// SetBorderTopStyle 设置边框样式：上
func (r *Cell) SetBorderTopStyle(borderTopStyle int, condition bool) *Cell {
	if condition {
		r.borderTopStyle = borderTopStyle
	}
	return r
}

// SetBorderTopStyleFunc 设置边框样式：上 函数
func (r *Cell) SetBorderTopStyleFunc(condition func() (int, bool)) *Cell {
	if condition != nil {
		r.SetBorderTopStyle(condition())
	}
	return r
}

// SetBorderBottomRgb 设置边框颜色：下
func (r *Cell) SetBorderBottomRgb(borderBottomRgb string, condition bool) *Cell {
	if condition {
		r.borderBottomRgb = borderBottomRgb
	}
	return r
}

// SetBorderBottomRbgFunc 设置边框颜色：下 函数
func (r *Cell) SetBorderBottomRbgFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		r.SetBorderBottomRgb(condition())
	}
	return r
}

// SetBorderBottomStyle 设置边框样式：下
func (r *Cell) SetBorderBottomStyle(borderBottomStyle int, condition bool) *Cell {
	if condition {
		r.borderBottomStyle = borderBottomStyle
	}
	return r
}

// SetBorderBottomStyleFunc 设置边框样式：下 函数
func (r *Cell) SetBorderBottomStyleFunc(condition func() (int, bool)) *Cell {
	if condition != nil {
		r.SetBorderBottomStyle(condition())
	}
	return r
}

// SetBorderLeftRgb 设置边框颜色：左
func (r *Cell) SetBorderLeftRgb(borderLeftRgb string, condition bool) *Cell {
	if condition {
		r.borderLeftRgb = borderLeftRgb
	}
	return r
}

// SetBorderLeftRbgFunc 设置边框颜色：左 函数
func (r *Cell) SetBorderLeftRbgFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		r.SetBorderLeftRgb(condition())
	}
	return r
}

// SetBorderLeftStyle 设置边框样式：左
func (r *Cell) SetBorderLeftStyle(borderLeftStyle int, condition bool) *Cell {
	if condition {
		r.borderLeftStyle = borderLeftStyle
	}
	return r
}

// SetBorderLeftStyleFunc 设置边框样式：左 函数
func (r *Cell) SetBorderLeftStyleFunc(condition func() (int, bool)) *Cell {
	if condition != nil {
		r.SetBorderLeftStyle(condition())
	}
	return r
}

// SetBorderRightRgb 设置边框颜色：右
func (r *Cell) SetBorderRightRgb(borderRightRgb string, condition bool) *Cell {
	if condition {
		r.borderRightRgb = borderRightRgb
	}
	return r
}

// SetBorderRightRbgFunc 设置边框颜色：右 函数
func (r *Cell) SetBorderRightRbgFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		r.SetBorderRightRgb(condition())
	}
	return r
}

// SetBorderRightStyle 设置边框样式：右
func (r *Cell) SetBorderRightStyle(borderRightStyle int, condition bool) *Cell {
	if condition {
		r.borderRightStyle = borderRightStyle
	}
	return r
}

// SetBorderRightStyleFunc 设置边框样式：右 函数
func (r *Cell) SetBorderRightStyleFunc(condition func() (int, bool)) *Cell {
	if condition != nil {
		r.SetBorderRightStyle(condition())
	}
	return r
}

// SetBorderDiagonalUpRgb 设置边框颜色：对角线上
func (r *Cell) SetBorderDiagonalUpRgb(borderDiagonalUpRgb string, condition bool) *Cell {
	if condition {
		r.borderDiagonalUpRgb = borderDiagonalUpRgb
	}
	return r
}

// SetBorderDiagonalUpRbgFunc 设置边框颜色：对角线上 函数
func (r *Cell) SetBorderDiagonalUpRbgFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		r.SetBorderDiagonalUpRgb(condition())
	}
	return r
}

// SetBorderDiagonalUpStyle 设置边框样式：对角线上
func (r *Cell) SetBorderDiagonalUpStyle(borderDiagonalUpStyle int, condition bool) *Cell {
	if condition {
		r.borderDiagonalUpStyle = borderDiagonalUpStyle
	}
	return r
}

// SetBorderDiagonalUpStyleFunc 设置边框样式：对角线上 函数
func (r *Cell) SetBorderDiagonalUpStyleFunc(condition func() (int, bool)) *Cell {
	if condition != nil {
		r.SetBorderDiagonalUpStyle(condition())
	}
	return r
}

// SetBorderDiagonalDownRgb 设置边框颜色：对角线下
func (r *Cell) SetBorderDiagonalDownRgb(borderDiagonalDownRgb string, condition bool) *Cell {
	if condition {
		r.borderDiagonalDownRgb = borderDiagonalDownRgb
	}
	return r
}

// SetBorderDiagonalDownRbgFunc 设置边框颜色：对角线下 函数
func (r *Cell) SetBorderDiagonalDownRbgFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		r.SetBorderDiagonalDownRgb(condition())
	}
	return r
}

// SetBorderDiagonalDownStyle 设置边框样式：对角线下
func (r *Cell) SetBorderDiagonalDownStyle(borderDiagonalDownStyle int, condition bool) *Cell {
	if condition {
		r.borderDiagonalDownStyle = borderDiagonalDownStyle
	}
	return r
}

// SetBorderDiagonalDownStyleFunc 设置边框样式：对角线下 函数
func (r *Cell) SetBorderDiagonalDownStyleFunc(condition func() (int, bool)) *Cell {
	if condition != nil {
		r.SetBorderDiagonalDownStyle(condition())
	}
	return r
}

// GetFontRgb 获取字体颜色
func (r *Cell) GetFontRgb() string {
	return r.fontRgb
}

// SetFontRgb 设置字体颜色
func (r *Cell) SetFontRgb(fontRgb string, condition bool) *Cell {
	if condition {
		r.fontRgb = fontRgb
	}
	return r
}

// SetFontRgbFunc 设置字体颜色：函数
func (r *Cell) SetFontRgbFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		r.SetFontRgb(condition())
	}
	return r
}

// GetPatternRgb 获取填充色
func (r *Cell) GetPatternRgb() string {
	return r.patternRgb
}

// SetPatternRgb 设置填充色
func (r *Cell) SetPatternRgb(patternRgb string, condition bool) *Cell {
	if condition {
		r.patternRgb = patternRgb
	}
	return r
}

// SetPatternRgbFunc 设置填充色：函数
func (r *Cell) SetPatternRgbFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		r.SetPatternRgb(condition())
	}
	return r
}

// GetFontBold 获取字体粗体
func (r *Cell) GetFontBold() bool {
	return r.fontBold
}

// SetFontBold 设置字体粗体
func (r *Cell) SetFontBold(fontBold bool, condition bool) *Cell {
	if condition {
		r.fontBold = fontBold
	}
	return r
}

// SetFontBoldFunc 设置字体粗体：函数
func (r *Cell) SetFontBoldFunc(condition func() (bool, bool)) *Cell {
	if condition != nil {
		r.SetFontBold(condition())
	}
	return r
}

// GetFontItalic 获取字体斜体
func (r *Cell) GetFontItalic() bool {
	return r.fontItalic
}

// SetFontItalic 设置字体斜体
func (r *Cell) SetFontItalic(fontItalic bool, condition bool) *Cell {
	if condition {
		r.fontItalic = fontItalic
	}
	return r
}

// SetFontItalicFunc 设置字体斜体：函数
func (r *Cell) SetFontItalicFunc(condition func() (bool, bool)) *Cell {
	if condition != nil {
		r.SetFontItalic(condition())
	}
	return r
}

// GetFontFamily 获取字体
func (r *Cell) GetFontFamily() string {
	return r.fontFamily
}

// SetFontFamily 设置字体
func (r *Cell) SetFontFamily(fontFamily string, condition bool) *Cell {
	if condition {
		r.fontFamily = fontFamily
	}
	return r
}

// SetFontFamilyFunc 设置字体：函数
func (r *Cell) SetFontFamilyFunc(condition func() (string, bool)) *Cell {
	if condition != nil {
		r.SetFontFamily(condition())
	}
	return r
}

// GetFontSize 获取字体字号
func (r *Cell) GetFontSize() float64 {
	return r.fontSize
}

// SetFontSize 设置字体字号
func (r *Cell) SetFontSize(fontSize float64, condition bool) *Cell {
	if condition {
		r.fontSize = fontSize
	}
	return r
}

// SetFontSizeFunc 设置字体字号：函数
func (r *Cell) SetFontSizeFunc(condition func() (float64, bool)) *Cell {
	if condition != nil {
		r.SetFontSize(condition())
	}
	return r
}

// Init 初始化
func (r *Cell) Init(content any) *Cell {
	r.content = content
	return r
}

// GetContent 获取内容
func (r *Cell) GetContent() any {
	return r.content
}

// SetContent 设置内容
func (r *Cell) SetContent(content any) *Cell {
	r.content = content
	return r
}

// GetCoordinate 获取单元格坐标
func (r *Cell) GetCoordinate() string {
	return r.coordinate
}

// SetCoordinate 设置单元格坐标
func (r *Cell) SetCoordinate(coordinate string) *Cell {
	r.coordinate = coordinate
	return r
}

// GetContentType 获取单元格类型
func (r *Cell) GetContentType() CellContentType {
	return r.contentType
}

// SetContentType 设置单元格类型
func (r *Cell) SetContentType(contentType CellContentType) *Cell {
	r.contentType = contentType
	return r
}

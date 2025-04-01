package execl

import (
	"fmt"
	"reflect"

	"github.com/charmbracelet/log"
	"github.com/xuri/excelize/v2"
)

type ExportAttr struct {
	Prop      string                          `json:"prop"`
	Label     string                          `json:"label"`
	Width     int                             `json:"width"`
	Align     string                          `json:"align"`
	Formatter func(v interface{}) interface{} `json:"formatter"`
}

type IExport interface {
	ExportAttrs() []ExportAttr
	Export(list any) (*excelize.File, error)
}

type Export struct {
	Fields []ExportAttr `json:"fields"`
}

var _ IExport = (*Export)(nil)

func (e *Export) ExportAttrs() []ExportAttr {
	return e.Fields
}

func (e *Export) Export(list any) (workbook *excelize.File, err error) {
	workbook = excelize.NewFile()
	_, err = workbook.NewSheet("Sheet1")

	if err != nil {
		log.Error(err)
		return nil, err
	}

	fields := e.ExportAttrs()
	AZ := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i := 0; i < len(fields); i++ {
		field := fields[i]

		if field.Prop == "" {
			continue
		}

		label := field.Label

		if label == "" {
			label = field.Prop
		}
		workbook.SetCellValue("Sheet1", fmt.Sprintf("%s1", AZ[i]), label)
		// set column width
		width := field.Width

		if width == 0 {
			width = 8
		}
		workbook.SetColWidth("Sheet1", AZ[i], AZ[i], float64(width))

		// set align
		style, err := workbook.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: field.Align,
			},
			Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#f2f2f2"}},
			Border: []excelize.Border{
				{Type: "left", Color: "#000000", Style: 1},
				{Type: "top", Color: "#000000", Style: 1},
				{Type: "right", Color: "#000000", Style: 1},
				{Type: "bottom", Color: "#000000", Style: 1},
			},
		})
		if err != nil {
			log.Error(err)
		}
		workbook.SetCellStyle("Sheet1", fmt.Sprintf("%s1", AZ[i]), fmt.Sprintf("%s1", AZ[i]), style)
	}

	listRef := reflect.ValueOf(list)

	if listRef.Kind() == reflect.Ptr {
		listRef = listRef.Elem()
	}

	if listRef.Kind() != reflect.Slice {
		return nil, fmt.Errorf("list must be a slice")
	}

	length := listRef.Len()

	for i := 0; i < length; i++ {
		el := listRef.Index(i)

		for j := 0; j < len(fields); j++ {
			field := fields[j]
			if field.Prop == "" {
				continue
			}

			filed := el.FieldByName(field.Prop)

			if !filed.IsValid() {
				continue
			}

			value := filed.Interface()
			formatter := field.Formatter

			if formatter != nil {
				value = formatter(value)
			}

			if value == nil {
				value = "-"
			}

			workbook.SetCellValue("Sheet1", fmt.Sprintf("%s%d", AZ[j], i+2), value)
			style, err := workbook.NewStyle(&excelize.Style{
				Alignment: &excelize.Alignment{
					Horizontal: field.Align,
				},
			})
			if err != nil {
				log.Error(err)
			}
			workbook.SetCellStyle("Sheet1", fmt.Sprintf("%s%d", AZ[j], i+2), fmt.Sprintf("%s%d", AZ[j], i+2), style)
		}
	}

	return workbook, nil
}

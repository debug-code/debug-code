package fileTools

import (
	"git.xrewin.com/go/beego"
	"github.com/tealeg/xlsx"
)

func OpenXlsxWithByte(bs []byte) (*xlsx.File, error) {
	xlf, err := xlsx.OpenBinary(bs)
	if err != nil {
		beego.Error(err)
		return xlf, err
	}
	//beego.Debug("xlf", xlf)

	return xlf, nil
}
func GetDataWithByte(bs []byte) ([][]string, error) {
	xlf, err := OpenXlsxWithByte(bs)
	if err != nil {
		return [][]string{}, err
	}

	cells := len(xlf.Sheets[0].Rows[0].Cells)
	data := [][]string{}
	for _, row := range xlf.Sheets[0].Rows {
		temp := []string{}
		index := 0
		for k, v := range row.Cells {
			index = k
			if (k + 1) > cells {
				break
			}
			//beego.Debug(k, v.String())
			temp = append(temp, v.String())
		}
		for i := 1; i < cells-index; i++ {
			temp = append(temp, "")
		}
		index = 0
		data = append(data, temp)
	}

	return data, nil
}

func OpenXlsx(path string) *xlsx.File {

	xlf, err := xlsx.OpenFile(path)
	if err != nil {
		beego.Error(err)
	}
	//beego.Debug("xlf", xlf)

	return xlf
}
func GetDatafrom(path string) [][]string {
	xlf := OpenXlsx(path)

	cells := len(xlf.Sheets[0].Rows[0].Cells)
	data := [][]string{}
	for _, row := range xlf.Sheets[0].Rows {
		temp := []string{}
		index := 0
		for k, v := range row.Cells {
			index = k
			if (k + 1) > cells {
				break
			}
			//beego.Debug(k, v.String())
			temp = append(temp, v.String())
		}
		for i := 1; i < cells-index; i++ {
			temp = append(temp, "")
		}
		index = 0
		data = append(data, temp)
	}

	return data
}

package utils

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
)

type Excel struct {

}

type ExcelUtilMgr struct {
}

var excelUtilMgr *ExcelUtilMgr = nil

func GetExcelUtilMgr() *ExcelUtilMgr {
	if excelUtilMgr == nil {
		excelUtilMgr = new(ExcelUtilMgr)
	}
	return excelUtilMgr
}

func (this* ExcelUtilMgr) ReadExcel(fileName string)[][]string  {

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		return [][]string{}
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 获取工作表中指定单元格的值
	//cell, err := f.GetCellValue("Sheet1", "B2")
	//if err != nil {
	//	fmt.Println(err)
	//	return [][]string{}
	//}
	//fmt.Println(cell)
	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return [][]string{}
	}
	//for _, row := range rows {
	//	//if i == 0 {
	//	//	continue
	//	//}
	//	for _, colCell := range row {
	//		fmt.Print(colCell, "\t")
	//	}
	//	fmt.Println()
	//}
	return rows

}


func (this* ExcelUtilMgr) LoadExcel(fileName string,SlicePtr interface{}) {
	excelFile := "excel/" + fileName + ".xlsx"
	excelData := this.ReadExcel(excelFile)
	if len(excelData) <= 1 {
		fmt.Println("len(csvData) <= 1, filename:", fileName)
		os.Exit(1)
		return
	}

	//LogDebug("Read File:", fileName, len(csvData))

	err := GetCsvUtilMgr().ParseDataSimple(excelData, SlicePtr, fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println("SlicePtr:",SlicePtr)
}
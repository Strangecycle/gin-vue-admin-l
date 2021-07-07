package service

import (
	"errors"
	"fmt"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"strconv"
)

func ParseInfoList2Excel(infoList []model.SysBaseMenu, filePath string) (err error) {
	excel := excelize.NewFile()
	// 设置第一行的数据（表头）
	err = excel.SetSheetRow("Sheet1", "A1", &[]string{"ID", "路由Name", "路由Path", "是否隐藏", "父节点", "排序", "文件名称"})
	for i, menu := range infoList {
		axis := fmt.Sprintf("A%d", i+2)
		err = excel.SetSheetRow("Sheet1", axis, &[]interface{}{
			menu.ID,
			menu.Name,
			menu.Path,
			menu.Hidden,
			menu.ParentId,
			menu.Sort,
			menu.Component,
		})
	}
	err = excel.SaveAs(filePath)
	return nil
}

func ParseExcel2InfoList() (menus []model.SysBaseMenu, err error) {
	skipHeader := true
	fixedHeader := []string{"ID", "路由Name", "路由Path", "是否隐藏", "父节点", "排序", "文件名称"}
	file, err := excelize.OpenFile(global.GVA_CONFIG.Excel.Dir + "ExcelImport.xlsx")
	if err != nil {
		return nil, err
	}

	menus = make([]model.SysBaseMenu, 0)
	rows, err := file.Rows("Sheet1")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		if skipHeader {
			// 对比表头
			if !compareStrSlice(row, fixedHeader) {
				return nil, errors.New("Excel 格式错误")
			}
			skipHeader = false
			continue
		}
		// 如果这一行与表头长度不一致，跳出当前这一行的执行
		if len(row) != len(fixedHeader) {
			continue
		}
		id, _ := strconv.Atoi(row[0])
		hidden, _ := strconv.ParseBool(row[3])
		sort, _ := strconv.Atoi(row[5])
		menu := model.SysBaseMenu{
			GVA_MODEL: global.GVA_MODEL{
				ID: uint(id),
			},
			Name:      row[1],
			Path:      row[2],
			Hidden:    hidden,
			ParentId:  row[4],
			Sort:      sort,
			Component: row[6],
		}
		menus = append(menus, menu)
	}
	return menus, nil
}

func compareStrSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	if (a == nil) != (b == nil) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

package reader

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"runtime/debug"
)

func WriteXlsx(path string, file *excelize.File) error {
	if err := os.MkdirAll(filepath.Dir(path), 644); err != nil {
		return err
	}

	return file.SaveAs(path)
}

func WriteXls(path string, file *xlsx.File) error {
	if err := os.MkdirAll(filepath.Dir(path), 644); err != nil {
		return err
	}

	return file.Save(path)
}

// 创建Office应用实例（优先WPS，其次Excel）
func createOfficeApplication() (*ole.IDispatch, error) {
	// 尝试创建WPS应用
	unknown, err := oleutil.CreateObject("wps.Application")
	if err == nil {
		fmt.Println("已连接到WPS")
		wps, _ := unknown.QueryInterface(ole.IID_IDispatch)
		return wps, nil
	}

	// 尝试创建ET（WPS表格）应用
	unknown, err = oleutil.CreateObject("et.Application")
	if err == nil {
		fmt.Println("已连接到WPS表格")
		et, _ := unknown.QueryInterface(ole.IID_IDispatch)
		return et, nil
	}

	// 尝试创建Excel应用
	unknown, err = oleutil.CreateObject("Excel.Application")
	if err != nil {
		return nil, fmt.Errorf("未找到WPS或Excel: %v", err)
	}

	fmt.Println("已连接到Excel")
	excel, _ := unknown.QueryInterface(ole.IID_IDispatch)
	return excel, nil
}

func ExcelXlsx2Xls(xlsxPath string, xlsPath string) error {
	//转化为绝对路径
	xlsxPath, err := filepath.Abs(xlsxPath)
	if err != nil {
		return err
	}

	xlsPath, err = filepath.Abs(xlsPath)
	if err != nil {
		return err
	}

	// 检查文件是否存在
	if _, err := os.Stat(xlsxPath); os.IsNotExist(err) {
		return fmt.Errorf("%v 文件不存在", xlsxPath)
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("捕获到panic:", r, string(debug.Stack())) // 输出panic传递的值
		}
	}()

	// 初始化COM库
	err = ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	if err != nil {
		return err
	}
	defer ole.CoUninitialize()

	// 尝试创建WPS/Excel应用实例
	app, err := createOfficeApplication()
	if err != nil {
		return fmt.Errorf("无法创建WPS或Excel应用: %v", err)
	}
	defer app.Release()

	// 设置不显示警告对话框
	oleutil.PutProperty(app, "DisplayAlerts", false)
	defer func() {
		// 恢复设置
		oleutil.PutProperty(app, "DisplayAlerts", true)
	}()

	// 打开.xlsx文件
	workbooks := oleutil.MustGetProperty(app, "Workbooks").ToIDispatch()
	defer workbooks.Release()

	workbook := oleutil.MustCallMethod(workbooks, "Open", xlsxPath).ToIDispatch()
	defer workbook.Release()

	// 保存为.xls格式（Excel 97-2003工作簿格式）
	const xlExcel8 = 56 // 对应Excel 97-2003格式的常量值
	oleutil.MustCallMethod(workbook, "SaveAs", xlsPath, xlExcel8, "", "", false, false)

	// 关闭工作簿
	oleutil.MustCallMethod(workbook, "Close")

	// 退出应用
	oleutil.MustCallMethod(app, "Quit")

	fmt.Printf("成功将 %s 转换为 %s\n", filepath.Base(xlsxPath), filepath.Base(xlsPath))

	return nil
}

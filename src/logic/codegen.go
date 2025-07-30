package logic

import (
	"Meta2XlsGen/src/cmd"
	"Meta2XlsGen/src/reader"
	"Meta2XlsGen/src/types/typedef"
	"Meta2XlsGen/src/utils"
	"bytes"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path"
	"text/template"
	"unicode/utf8"
)

type GenCodeLogic struct {
	args *cmd.CmdArgs
}

func NewGenCodeLogic(args *cmd.CmdArgs) *GenCodeLogic {
	return &GenCodeLogic{
		args: args,
	}
}

func (l *GenCodeLogic) Run() error {
	//构造模版参数
	tplArgs := NewTemplateArgs(l.args)
	err := tplArgs.genArgs()
	if err != nil {
		return err
	}

	//生成代码文件
	err = l.genCode("ConvertConfig.tpl", tplArgs, fmt.Sprintf("%sConvMsg.txt", "Tmp"), true)
	if err != nil {
		return err
	}

	if !l.args.OnlyConvMsg {
		if !l.args.NativeMode {
			xlsxFile, err := l.genXlsx(tplArgs)
			if err != nil {
				return err
			}

			pth := path.Join(l.args.ExcelPath, l.args.ExcelFile)
			tmpPath := utils.RepleacePathExt(pth, "Tmp.xlsx")
			err = reader.WriteXlsx(tmpPath, xlsxFile)
			if err != nil {
				return err
			}
			log.Printf("gen xlsx file:%v\n", tmpPath)

			log.Println("begin convert")
			err = reader.ExcelXlsx2Xls(tmpPath, pth)
			if err != nil {
				return err
			}

			err = os.Remove(tmpPath)
			if err != nil {
				return err
			}
		} else {
			xlsFile, err := l.genXls(tplArgs)
			if err != nil {
				return err
			}

			pth := path.Join(l.args.ExcelPath, l.args.ExcelFile)
			err = reader.WriteXls(pth, xlsFile)
			if err != nil {
				return err
			}
			log.Printf("gen xls file:%v\n", pth)
		}
	}

	return nil
}

func (l *GenCodeLogic) genCode(tplName string, tplArgs *TemplateArgs, outName string, overwrite bool) error {
	fileName := path.Join(l.args.TmpPath, outName)
	if !overwrite {
		_, err := os.Stat(fileName)
		if err == nil { //文件存在
			return nil //跳过
		}
	}

	//读取模版文件
	tplData, err := os.ReadFile(path.Join(utils.GetExecutablePath(), "templates", tplName))
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New("codeTemplate").Funcs(tplArgs.GetFuncMap()).Parse(string(tplData)))

	// 创建一个缓冲区来保存渲染后的内容
	var rendered bytes.Buffer

	// 渲染模板
	err = tmpl.Execute(&rendered, tplArgs)
	if err != nil {
		return err
	}

	outBytes := rendered.Bytes()
	//if l.args.Encoding == "GBK" {
	//	encoder := simplifiedchinese.GBK.NewEncoder()
	//	reader := transform.NewReader(bytes.NewReader(outBytes), encoder)
	//	outBytes, err = io.ReadAll(reader)
	//	if err != nil {
	//		return err
	//	}
	//}

	// 写入文件
	err = os.WriteFile(fileName, outBytes, 0644)
	if err != nil {
		return err
	}
	log.Printf("gen tmp file:%v, please read it\n", fileName)
	return nil
}

type content struct {
	title    string
	comment  string
	firstRow interface{}
}

func (l *GenCodeLogic) genXls(tplArgs *TemplateArgs) (*xlsx.File, error) {
	file := xlsx.NewFile()

	font := xlsx.Font{
		Name: "Microsoft YaHei UI", // 字体名称
		Size: 11,                   // 字号
		Bold: true,                 // 是否加粗
	}

	for _, s := range tplArgs.XmlFile.Defines {
		if !s.TagOption.Export {
			continue
		}

		sheet, err := file.AddSheet(s.Desc)
		if err != nil {
			return nil, err
		}

		cols := make([]*content, 0)
		for _, field := range s.Field {
			err = l.genCell(tplArgs, &cols, field, "", "")
			if err != nil {
				return nil, err
			}
		}

		row := sheet.AddRow()
		for i, val := range cols {
			cell := row.AddCell()
			cell.SetString(val.title)

			//样式
			cell.GetStyle().Font = font
			cell.GetStyle().Alignment.Horizontal = "center"
			cell.GetStyle().ApplyAlignment = true
			width := float64(len(val.title)) * 1.1
			if utf8.RuneCountInString(val.title) < len(val.title) { //带中文
				width -= float64(len(val.title)-utf8.RuneCountInString(val.title)) * 0.65
			}
			sheet.Col(i).Width = utils.Max(width, 5.5)
		}

		//数据
		count := 1
		if !s.TagOption.IsSingleLine {
			count = 3
		}
		for rowCount := 0; rowCount < count; rowCount++ {
			row = sheet.AddRow()
			for _, val := range cols {
				cell := row.AddCell()
				cell.SetValue(val.firstRow)

				//样式
				cell.GetStyle().Font.Name = "Microsoft YaHei UI"
				cell.GetStyle().Font.Size = 10
				cell.GetStyle().Alignment.Horizontal = "center"
				cell.GetStyle().Alignment.WrapText = true
				cell.GetStyle().ApplyAlignment = true
			}
		}
	}

	return file, nil
}

func (l *GenCodeLogic) genXlsx(tplArgs *TemplateArgs) (*excelize.File, error) {
	file := excelize.NewFile()

	//样式
	fontTitle := &excelize.Font{Family: "Microsoft YaHei UI", Size: 11, Bold: true}
	alignmentTitle := &excelize.Alignment{Horizontal: "center", Vertical: "center"}
	styleTitle, err := file.NewStyle(&excelize.Style{Font: fontTitle, Alignment: alignmentTitle})
	if err != nil {
		return nil, err
	}

	fontBody := &excelize.Font{Family: "Microsoft YaHei UI", Size: 10}
	alignmentBody := &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true}
	styleBody, err := file.NewStyle(&excelize.Style{Font: fontBody, Alignment: alignmentBody})
	if err != nil {
		return nil, err
	}

	for _, s := range tplArgs.XmlFile.Defines {
		if !s.TagOption.Export {
			continue
		}

		sheetName := s.Desc
		_, err := file.NewSheet(sheetName)
		if err != nil {
			return nil, err
		}

		cols := make([]*content, 0)

		for _, field := range s.Field {
			err = l.genCell(tplArgs, &cols, field, "", "")
			if err != nil {
				return nil, err
			}
		}

		for i, val := range cols {
			col, err := excelize.ColumnNumberToName(i + 1)
			if err != nil {
				return nil, err
			}

			cell := fmt.Sprintf("%s%d", col, 1)
			if err = file.SetCellValue(sheetName, cell, val.title); err != nil {
				return nil, err
			}

			//列宽度
			width := float64(len(val.title)) * 1.5
			if utf8.RuneCountInString(val.title) < len(val.title) { //带中文
				width -= float64(len(val.title)-utf8.RuneCountInString(val.title)) * 0.8
			}
			if err = file.SetColWidth(sheetName, col, col, utils.Max(width, 10)); err != nil {
				return nil, err
			}

			//列样式
			if err = file.SetColStyle(sheetName, col, styleBody); err != nil {
				return nil, err
			}

			//标题样式
			if err = file.SetCellStyle(sheetName, cell, cell, styleTitle); err != nil {
				return nil, err
			}

			//批注
			if len(val.comment) > 0 {
				if err = file.AddComment(sheetName, excelize.Comment{
					Author: "Meta2XlsGen",
					Cell:   cell,
					Text:   val.comment,
				}); err != nil {
					return nil, err
				}
			}
		}

		//数据
		count := 1
		if !s.TagOption.IsSingleLine {
			count = 3
		}
		for row := 0; row < count; row++ {
			for i, val := range cols {
				col, err := excelize.ColumnNumberToName(i + 1)
				if err != nil {
					return nil, err
				}
				cell := fmt.Sprintf("%s%d", col, row+2)
				if err = file.SetCellValue(sheetName, cell, val.firstRow); err != nil {
					return nil, err
				}
			}
		}
	}

	err = file.DeleteSheet("Sheet1")
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (l *GenCodeLogic) genCell(tplArgs *TemplateArgs, cols *[]*content, field *typedef.FieldInfo, prefix, prefixComent string) error {
	comment := prefixComent
	if len(comment) > 0 || len(field.Desc) > 0 {
		comment = fmt.Sprintf("%v\n%v", prefixComent, field.Desc)
	}

	if field.Type == typedef.FTStruct {
		s := tplArgs.XmlFile.FindDefineByName(field.TypeName)

		if s == nil {
			return fmt.Errorf("not found type %v", field.TypeName)
		}

		if field.Count > 0 {
			for i := 1; i <= field.Count; i++ {
				for _, f := range s.Field {
					err := l.genCell(tplArgs, cols, f,
						fmt.Sprintf("%s%s%d", prefix, field.CName, i), comment)
					if err != nil {
						return err
					}
				}
			}
		} else {
			for _, f := range s.Field {
				err := l.genCell(tplArgs, cols, f, prefix+field.CName, comment)
				if err != nil {
					return err
				}
			}
		}
	} else {
		if field.Count > 0 {
			return fmt.Errorf("not support array execpt struct, field:%v", field.Name)
		}
		data := &content{
			title:    prefix + field.CName,
			comment:  comment,
			firstRow: field.DefaultValue(),
		}
		*cols = append(*cols, data)
	}

	return nil
}

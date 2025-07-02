package logic

import (
	"Meta2XlsGen/src/cmd"
	"Meta2XlsGen/src/reader"
	"Meta2XlsGen/src/types/typedef"
	"Meta2XlsGen/src/utils"
	"bytes"
	"fmt"
	"github.com/tealeg/xlsx"
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
	err = l.genCode("ConvertConfig.tpl", tplArgs, fmt.Sprintf("%sConvMsg.txt", l.args.Name), true)
	if err != nil {
		return err
	}

	xlsFile, err := l.genXls(tplArgs)
	if err != nil {
		return err
	}

	pth := path.Join(l.args.OutPath, l.args.ExcelFile)
	err = reader.WriteXls(pth, xlsFile)
	if err != nil {
		return err
	}
	log.Printf("gen xls file:%v\n", pth)

	return nil
}

func (l *GenCodeLogic) genCode(tplName string, tplArgs *TemplateArgs, outName string, overwrite bool) error {
	fileName := path.Join(l.args.OutPath, outName)
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
	log.Printf("gen xml file:%v\n", fileName)
	return nil
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

		colsTitle := make([]string, 0)
		colsDefaultData := make([]interface{}, 0)

		for _, field := range s.Field {
			err = l.genCell(tplArgs, &colsTitle, &colsDefaultData, field, "")
			if err != nil {
				return nil, err
			}
		}

		row := sheet.AddRow()
		for i, col := range colsTitle {
			cell := row.AddCell()
			cell.SetString(col)

			//样式
			cell.GetStyle().Font = font
			cell.GetStyle().Alignment.Horizontal = "center"
			cell.GetStyle().ApplyAlignment = true
			width := float64(len(col)) * 1.1
			if utf8.RuneCountInString(col) < len(col) { //带中文
				width -= float64(len(col)-utf8.RuneCountInString(col)) * 0.65
			}
			sheet.Col(i).Width = utils.Max(width, 5.5)

		}

		row = sheet.AddRow()
		for _, col := range colsDefaultData {
			cell := row.AddCell()
			cell.SetValue(col)

			//样式
			cell.GetStyle().Font.Name = "Microsoft YaHei UI"
			cell.GetStyle().Font.Size = 10
			cell.GetStyle().Alignment.Horizontal = "center"
			cell.GetStyle().Alignment.WrapText = true
			cell.GetStyle().ApplyAlignment = true
		}
	}

	return file, nil
}

func (l *GenCodeLogic) genCell(tplArgs *TemplateArgs, cols *[]string, defaultData *[]interface{}, field *typedef.FieldInfo, prefix string) error {
	if field.Type == typedef.FTStruct {
		var s *typedef.StructInfo
		for _, v := range tplArgs.XmlFile.Defines {
			if v.Name == field.TypeName {
				s = v
				break
			}
		}

		if s == nil {
			return fmt.Errorf("not found type %v", field.TypeName)
		}

		if field.Count > 0 {
			for i := 1; i <= field.Count; i++ {
				for _, f := range s.Field {
					err := l.genCell(tplArgs, cols, defaultData, f, fmt.Sprintf("%s%s%d", prefix, field.CName, i))
					if err != nil {
						return err
					}
				}
			}
		} else {
			for _, f := range s.Field {
				err := l.genCell(tplArgs, cols, defaultData, f, prefix+field.CName)
				if err != nil {
					return err
				}
			}
		}
	} else {
		if field.Count > 0 {
			return fmt.Errorf("not support array execpt struct, field:%v", field.Name)
		}
		*cols = append(*cols, prefix+field.CName)
		*defaultData = append(*defaultData, field.DefaultValue())
	}

	return nil
}

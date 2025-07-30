package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path"
	"regexp"
	"slices"
)

var SupportEncoding []string = []string{
	"UTF-8",
	"GBK",
}

type CmdArgs struct {
	Version   string
	Name      string
	Author    string
	XmlPath   string
	Encoding  string
	LabelTag  map[string]string
	TmpPath   string
	ExcelPath string
	ExcelFile string

	NativeMode  bool
	OnlyConvMsg bool

	XmlModulePath string
	CodePath      string

	FilePath string

	SuccessRun bool
}

func (c *CmdArgs) parseCustomTypeLabel(strs []string) error {
	c.LabelTag = make(map[string]string)
	for _, labels := range strs {
		re := regexp.MustCompile(`^([\w.]+):<([^>]+)>$`)
		matches := re.FindAllStringSubmatch(labels, -1)
		if matches == nil {
			return fmt.Errorf("lable tag %v format error", labels)
		}
		for _, m := range matches {
			if len(m) == 3 {
				c.LabelTag[m[1]] = m[2]
			}
		}
	}
	return nil
}

const AppVersion = "v1.0.0"

func ParseCmdArgs() (*CmdArgs, error) {
	cmdArgs := &CmdArgs{
		Version: AppVersion,
	}
	cmd := &cobra.Command{
		Use: `tool for meta to xls
[]string类型的参数，支持多次输入，或者用,分隔，例:
  --ignore-label label1 --ignore-label label2 --ignore-label label3
  --ignore-label label1,label2 --ignore-label label3
`,
		Short: "meta配置转化工具",
		Long: `meta.xml文件，生成对应的xls文件

xml文件支持在标签上方添加注解tag，语法为<!--tag:"Key1:Value2,Key2:Value2"-->
目前支持的tag有：
	export: 是否导出，默认为true
	ignore: 默认为false,用于标记该entry是否需要再代码中生成（值字段生效）
	id: 默认为false,用于标记该entry是id字段（值字段生效）
	customType: 自定义类型的名称
	singleLine: 默认为false,用于标记该struct在xls中是否为单行数据，决定C++代码中是否为数组
`,
		Version: AppVersion,
		Example: `Meta2XlsGen --name TestAct file1 file2`,
		Args:    cobra.ArbitraryArgs,
		RunE: func(c *cobra.Command, args []string) error {
			var err error
			cmdArgs.Name, _ = c.Flags().GetString("name")
			cmdArgs.Author, _ = c.Flags().GetString("author")

			cmdArgs.XmlPath, _ = c.Flags().GetString("xmlPath")
			if err = os.MkdirAll(cmdArgs.XmlPath, 644); err != nil {
				return fmt.Errorf("invalid XmlPath err:%v", err.Error())
			}

			cmdArgs.TmpPath, _ = c.Flags().GetString("tmpPath")
			if err = os.MkdirAll(cmdArgs.TmpPath, 644); err != nil {
				return fmt.Errorf("invalid TmpPath err:%v", err.Error())
			}

			cmdArgs.CodePath, _ = c.Flags().GetString("codePath")

			cmdArgs.XmlModulePath = path.Join(cmdArgs.XmlPath, cmdArgs.Name)
			if err = os.MkdirAll(cmdArgs.XmlModulePath, 644); err != nil {
				return fmt.Errorf("invalid XmlModulePath err:%v", err.Error())
			}

			cmdArgs.ExcelPath, _ = c.Flags().GetString("excelPath")
			if err = os.MkdirAll(cmdArgs.ExcelPath, 644); err != nil {
				return fmt.Errorf("invalid ExcelPath err:%v", err.Error())
			}

			cmdArgs.ExcelFile, _ = c.Flags().GetString("excelFile")
			if path.Ext(cmdArgs.ExcelFile) != ".xls" {
				return fmt.Errorf("invalid excelFile path:%v", cmdArgs.ExcelFile)
			}

			cmdArgs.Encoding, _ = c.Flags().GetString("encoding")
			if !slices.Contains(SupportEncoding, cmdArgs.Encoding) {
				return fmt.Errorf("invalid encoding %v", cmdArgs.Encoding)
			}

			cmdArgs.NativeMode, _ = c.Flags().GetBool("nativeMode")
			cmdArgs.OnlyConvMsg, _ = c.Flags().GetBool("onlyConvMsg")

			LabelTag, _ := c.Flags().GetStringArray("label-tag")
			err = cmdArgs.parseCustomTypeLabel(LabelTag)
			if err != nil {
				return fmt.Errorf("parseCustomTypeLabel err:%v", err.Error())
			}

			if len(args) != 1 {
				return errors.New("文件路径数量不正确！")
			}

			cmdArgs.FilePath = args[0]

			cmdArgs.SuccessRun = true
			return nil
		},
	}

	cmd.Flags().String("name", "Tmp", "指定生成的配置名称，用于文件名、结构名等")
	cmd.Flags().String("author", "Meta2XlsGen", "生成作者")
	cmd.Flags().String("xmlPath", "./xls", "xml文件生成目录")
	cmd.Flags().String("tmpPath", "./", "tmp描述文件目录")
	cmd.Flags().String("codePath", "./", "代码文件生成目录，对应Xml2CodeGen out-path")
	cmd.Flags().String("excelPath", "./xls", "xls文件生成目录")
	cmd.Flags().String("excelFile", "测试目录/测试配置.xls", "xls文件路径")
	cmd.Flags().String("encoding", "GBK", "文件编码(UTF-8,GBK)")
	cmd.Flags().Bool("nativeMode", false, "是否原生模式，false:先生成xlsx文件，然后调用wps或excel软件转化为xls true:直接创建xls（不支持批注）")
	cmd.Flags().Bool("onlyConvMsg", false, "是否只生成转化信息，TmpConvMsg.txt")
	cmd.Flags().StringArray("label-tag", []string{}, "xml标签与tag的映射关系，用于省去xml文件中tag注释，label:`tag`")
	if err := cmd.Execute(); err != nil {
		return nil, err
	}

	return cmdArgs, nil
}

package logic

import (
	"Meta2XlsGen/src/cmd"
	"Meta2XlsGen/src/reader"
	"Meta2XlsGen/src/types/typedef"
	"Meta2XlsGen/src/utils"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type TemplateArgs struct {
	CmdArgs *cmd.CmdArgs

	Version    string
	CreateDate string

	XmlFile *FileInfo
}

type FileInfo struct {
	FileName string
	Macro    map[string]int
	Defines  []*typedef.StructInfo //结构变量，用于描述配置的属性结构
}

func NewTemplateArgs(c *cmd.CmdArgs) *TemplateArgs {
	return &TemplateArgs{
		CmdArgs:    c,
		Version:    c.Version,
		CreateDate: time.Now().Format(time.DateOnly),
	}
}

func (t *TemplateArgs) genArgs() error {
	//读取xml配置
	var err error
	t.XmlFile, err = t.parseXml(t.CmdArgs.FilePath)
	if err != nil {
		return err
	}

	labelNameMap := make(map[string]bool)
	for _, s := range t.XmlFile.Defines {
		//所有标签元素，名称不能重复
		if labelNameMap[s.Name] {
			return fmt.Errorf("has repeat label name:%v", s.Name)
		}
		labelNameMap[s.Name] = true

		for _, f := range s.Field {
			if c, ok := t.XmlFile.Macro[f.CountName]; ok {
				f.Count = c
			}
		}
	}

	return nil
}

func (t *TemplateArgs) GetFuncMap() template.FuncMap {
	return map[string]any{
		"toUpper": strings.ToUpper,
	}
}

func (t *TemplateArgs) parseXml(path string) (*FileInfo, error) {
	//读取文件
	info, err := reader.ReadXml(path)
	if err != nil {
		return nil, err
	}

	//xml信息解析
	xmlFile, err := t.parseMetaXml(info)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(path)
	xmlFile.FileName = name

	return xmlFile, nil
}

func (t *TemplateArgs) parseMetaXml(e *reader.Element) (*FileInfo, error) {
	if e.XMLName.Local != "metalib" {
		return nil, fmt.Errorf("root must be metalib label")
	}

	xmlFile := &FileInfo{Macro: make(map[string]int)}

	//遍历第一层定义
	for _, v := range e.Children {
		switch v.XMLName.Local {
		case "struct":
			ret, err := t.parseStruct(&v)
			if err != nil {
				return nil, err
			}
			xmlFile.Defines = append(xmlFile.Defines, ret)
		case "macro":
			key, val, err := t.parseMacro(&v)
			if err != nil {
				return nil, err
			}
			xmlFile.Macro[key] = val
		default:
			return nil, fmt.Errorf("meta child not support %v", v.XMLName.Local)
		}

	}
	return xmlFile, nil
}

func (t *TemplateArgs) parseMacro(e *reader.Element) (string, int, error) {
	name := utils.GetXmlAttr(e.Attrs, "name")
	if len(name) == 0 {
		return "", 0, fmt.Errorf("name is null")
	}

	return name, utils.GetXmlAttrInt(e.Attrs, "value"), nil
}

func (t *TemplateArgs) parseStruct(e *reader.Element) (*typedef.StructInfo, error) {
	name := utils.GetXmlAttr(e.Attrs, "name")
	if len(name) == 0 {
		return nil, fmt.Errorf("name is null")
	}

	//合并命令行中的tag信息
	if tagStr, ok := t.CmdArgs.LabelTag[name]; ok {
		op := typedef.NewTagOption(tagStr)
		e.TagOption.Merge(op)
	}

	s := &typedef.StructInfo{
		Name:      name,
		Desc:      utils.GetXmlAttr(e.Attrs, "desc"),
		Attr:      e.Attrs,
		Field:     make([]*typedef.FieldInfo, 0),
		TagOption: e.TagOption,
	}

	//属性成员
	for _, v := range e.Children {
		if v.XMLName.Local != "entry" {
			continue
		}

		f, err := t.parseEntry(&v)
		if err != nil {
			return nil, err
		}

		s.Field = append(s.Field, f)
	}

	return s, nil
}

func (t *TemplateArgs) parseEntry(e *reader.Element) (*typedef.FieldInfo, error) {
	name := utils.GetXmlAttr(e.Attrs, "name")
	if len(name) == 0 {
		return nil, fmt.Errorf("name is null")
	}

	tp := utils.GetXmlAttr(e.Attrs, "type")
	field := &typedef.FieldInfo{
		Name:      name,
		TypeName:  tp,
		CountName: utils.GetXmlAttr(e.Attrs, "count"),
		Count:     utils.GetXmlAttrInt(e.Attrs, "count"),
		Refer:     utils.GetXmlAttr(e.Attrs, "refer"),
		CName:     utils.GetXmlAttr(e.Attrs, "cname"),
		Desc:      utils.GetXmlAttr(e.Attrs, "desc"),
	}

	switch tp {
	case "int":
		field.Type = typedef.FTInt
	case "string":
		field.Type = typedef.FTString
	case "":
		return nil, fmt.Errorf("type is null")
	default:
		field.Type = typedef.FTStruct
	}

	return field, nil
}

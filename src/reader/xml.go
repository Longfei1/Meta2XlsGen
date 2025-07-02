package reader

import (
	"Meta2XlsGen/src/types/typedef"
	"encoding/xml"
	"fmt"
	"github.com/beevik/etree"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"regexp"
)

// 一开始使用的标准xml库解析，后面改为etree库
// 为了兼容性，保留了原始的结构定义
type Element struct {
	XMLName   xml.Name
	Attrs     []xml.Attr `xml:",any,attr"`
	Children  []Element  `xml:",any"`
	TagOption *typedef.TagOption
}

func ReadXml(path string) (*Element, error) {
	doc := etree.NewDocument()
	doc.ReadSettings.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "GBK" || charset == "gbk" {
			return transform.NewReader(input, simplifiedchinese.GBK.NewDecoder()), nil
		}
		return nil, fmt.Errorf("unsupported charset: %s", charset)
	}
	doc.ReadSettings.PreserveCData = false

	err := doc.ReadFromFile(path)
	if err != nil {
		return nil, err
	}

	root, err := parseDocument(doc)
	if err != nil {
		return nil, err
	}

	return root, nil
}

// 将etree.Element转换为Element，主要为处理特殊的tag注释
func parseDocument(doc *etree.Document) (*Element, error) {
	inRoot := doc.Root()
	return parseElement(inRoot)
}

func parseTagOption(e *etree.Element) *typedef.TagOption {
	if e == nil {
		return nil
	}

	var tagStr string

	//向前查找第一个合法的tag注释
	index := e.Index()
	for i := index - 1; i >= 0; i-- {
		t := e.Parent().Child[i]
		if c, ok := t.(*etree.Comment); ok {
			//解析tag
			re := regexp.MustCompile(`tag:"([^"]+)"`)
			matchStrs := re.FindStringSubmatch(c.Data)
			if len(matchStrs) != 2 {
				continue
			}

			tagStr = matchStrs[1]
			break
		}

		if cd, ok := t.(*etree.CharData); ok {
			if cd.IsWhitespace() {
				continue
			}
		}

		//其他类型的节点，终止遍历
		break
	}

	return typedef.NewTagOption(tagStr)
}

func parseElement(e *etree.Element) (*Element, error) {
	out := &Element{
		XMLName: xml.Name{
			Local: e.Tag,
		},
		TagOption: parseTagOption(e), //tag解析
	}

	//属性
	for _, v := range e.Attr {
		out.Attrs = append(out.Attrs, xml.Attr{
			Name: xml.Name{
				Local: v.Key,
			},
			Value: v.Value,
		})
	}

	//子节点
	for _, v := range e.ChildElements() {
		child, err := parseElement(v)
		if err != nil {
			return nil, err
		}
		out.Children = append(out.Children, *child)
	}

	return out, nil
}

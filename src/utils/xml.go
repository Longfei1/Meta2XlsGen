package utils

import (
	"encoding/xml"
	"strconv"
)

func GetXmlAttr(attrs []xml.Attr, name string) string {
	for _, v := range attrs {
		if v.Name.Local == name {
			return v.Value
		}
	}
	return ""
}

func GetXmlAttrInt(attrs []xml.Attr, name string) int {
	ret := GetXmlAttr(attrs, name)
	if len(ret) > 0 {
		n, _ := strconv.ParseInt(ret, 10, 32)
		return int(n)
	}
	return 0
}

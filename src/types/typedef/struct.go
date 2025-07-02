package typedef

import "encoding/xml"

type StructInfo struct {
	Name      string
	Desc      string
	Attr      []xml.Attr
	Field     []*FieldInfo
	TagOption *TagOption
}

func (s *StructInfo) AllFields() []*FieldInfo {
	return s.Field
}

func (s *StructInfo) FieldTypeCount(tp ...FieldType) int {
	cnt := 0
	for _, t := range tp {
		for _, v := range s.AllFields() {
			if v.Type == t {
				cnt++
			}
		}
	}
	return cnt
}

func (s *StructInfo) FieldByType(tp ...FieldType) []*FieldInfo {
	ret := make([]*FieldInfo, 0, len(s.Field))

	for _, t := range tp {
		for _, v := range s.AllFields() {
			if v.Type == t {
				ret = append(ret, v)
			}
		}
	}
	return ret
}

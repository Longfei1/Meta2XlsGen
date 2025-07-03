package typedef

import (
	"encoding/xml"
	"fmt"
)

type LabelTag struct {
	name string
	s    string
}

func NewLabelTag(name string) *LabelTag {
	return &LabelTag{
		name: name,
	}
}

func (l *LabelTag) Add(key string, val string) {
	if len(l.s) == 0 {
		l.s = fmt.Sprintf("%v:%v", key, val)
	} else {
		l.s = fmt.Sprintf("%v,%v:%v", l.s, key, val)
	}
}

func (l *LabelTag) IsEmpty() bool {
	return len(l.s) == 0
}

func (l *LabelTag) Value() string {
	return fmt.Sprintf("%v:<%v>", l.name, l.s)
}

type StructInfo struct {
	Name      string
	Desc      string
	Attr      []xml.Attr
	Field     []*FieldInfo
	TagOption *TagOption

	LabelTags []*LabelTag
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

func (s *StructInfo) FieldByName(name string) *FieldInfo {
	for _, v := range s.AllFields() {
		if v.Name == name {
			return v
		}
	}

	return nil
}

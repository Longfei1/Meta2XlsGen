package typedef

import (
	"encoding/xml"
	"fmt"
)

type LabelTag struct {
	name string
	k    []string
	kv   map[string]string
}

func NewLabelTag(name string) *LabelTag {
	return &LabelTag{
		name: name,
		k:    make([]string, 0),
		kv:   make(map[string]string),
	}
}

func (l *LabelTag) Add(key string, val string) {
	l.k = append(l.k, key)
	l.kv[key] = val
}

func (l *LabelTag) IsEmpty() bool {
	return len(l.k) == 0 || len(l.kv) == 0
}

func (l *LabelTag) Value() string {
	var s string
	for i, k := range l.k {
		if i == 0 {
			s = fmt.Sprintf("%v:%v", k, l.kv[k])
		} else {
			s = fmt.Sprintf("%v,%v:%v", s, k, l.kv[k])
		}
	}
	return fmt.Sprintf("%v:<%v>", l.name, s)
}

type StructInfo struct {
	Name      string
	Desc      string
	Attr      []xml.Attr
	Field     []*FieldInfo
	TagOption *TagOption

	IdNames     []string
	IgnoreAttr  []string
	FieldRemark []string
	FieldGetter []string
	SplitType   []string

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

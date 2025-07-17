package typedef

import (
	"strconv"
	"strings"
)

type TagKey string

const (
	TKExport      TagKey = "export"
	TKIgnore      TagKey = "ignore"
	TKId          TagKey = "id"
	TKCustomType  TagKey = "customType"
	TKSingleLine  TagKey = "singleLine"
	TKFieldGetter TagKey = "fieldGetter"
)

type TagInfo struct {
	Key   TagKey
	Value string
}

type TagOption struct {
	Export         bool
	IsId           bool
	IsIgnore       bool
	CustomTypeName string
	IsSingleLine   bool
	GetterName     string
}

func NewTagOption(tagStr string) *TagOption {
	//默认值
	op := &TagOption{
		Export: true,
	}

	if len(tagStr) > 0 {
		op.Parse(tagStr)
	}
	return op
}

func (t *TagOption) Parse(tagStr string) {
	tagStrs := strings.Split(tagStr, ",")
	for _, s := range tagStrs {
		if s == "" {
			continue
		}

		tagInfoStrs := strings.Split(s, ":")
		if len(tagInfoStrs) != 2 {
			continue
		}

		valueStr := tagInfoStrs[1]
		switch TagKey(tagInfoStrs[0]) {
		case TKExport:
			t.Export, _ = strconv.ParseBool(valueStr)
		case TKId:
			t.IsId, _ = strconv.ParseBool(valueStr)
		case TKIgnore:
			t.IsIgnore, _ = strconv.ParseBool(valueStr)
		case TKCustomType:
			t.CustomTypeName = valueStr
		case TKSingleLine:
			t.IsSingleLine, _ = strconv.ParseBool(valueStr)
		case TKFieldGetter:
			t.GetterName = valueStr
		}
	}
}

func (t *TagOption) Merge(src *TagOption) {
	if src == nil {
		return
	}

	if !src.Export {
		t.Export = src.Export
	}

	if !src.IsId {
		t.IsId = src.IsId
	}

	if src.IsIgnore {
		t.IsIgnore = src.IsIgnore
	}

	if len(src.CustomTypeName) > 0 {
		t.CustomTypeName = src.CustomTypeName
	}

	if len(src.GetterName) > 0 {
		t.GetterName = src.GetterName
	}
}

package typedef

type FieldType int

const (
	FTString FieldType = 1
	FTInt    FieldType = 2
	FTStruct FieldType = 3
)

type FieldInfo struct {
	Name      string
	Type      FieldType
	TypeName  string
	Count     int
	CountName string
	Refer     string
	CName     string
	Desc      string
}

func (f *FieldInfo) DefaultValue() interface{} {
	if f.Type == FTInt {
		return 0
	}
	return ""
}

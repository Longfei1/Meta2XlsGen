package typedef

type CustomType interface {
	TypeName() string
	FieldsInfo() []*FieldInfo //字段信息

	//模版代码相关
	LoadCode(varName string) string //加载数据代码

}

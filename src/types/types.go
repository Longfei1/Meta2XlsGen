package types

import "Meta2XlsGen/src/types/typedef"

var allTypes = make(map[string]typedef.CustomType)

func RegisterType(t typedef.CustomType) {
	allTypes[t.TypeName()] = t
}

func GetType(name string) typedef.CustomType {
	return allTypes[name]
}

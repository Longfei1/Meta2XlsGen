package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"path/filepath"
	"strings"
)

func ToCamelCase(s string) string {
	if len(s) == 0 {
		return s
	}

	// 按下划线分割字符串
	parts := strings.Split(s, "_")

	// 遍历每个部分，将首字母大写
	if len(parts) > 1 {
		caser := cases.Title(language.Und)
		for i, part := range parts {
			parts[i] = caser.String(part)
		}
		// 拼接成一个字符串
		return strings.Join(parts, "")
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func RepleacePathExt(path string, ext string) string {
	oldExt := filepath.Ext(path)                            //获取扩展名
	base := strings.TrimSuffix(filepath.Base(path), oldExt) // 去掉原扩展名
	dir := filepath.Dir(path)                               // 获取目录路径
	newPath := filepath.Join(dir, base+ext)                 // 拼接新路径
	return newPath
}

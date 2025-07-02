package reader

import (
	"github.com/tealeg/xlsx"
	"os"
	"path/filepath"
)

func WriteXls(path string, file *xlsx.File) error {
	if err := os.MkdirAll(filepath.Dir(path), 644); err != nil {
		return err
	}

	return file.Save(path)
}

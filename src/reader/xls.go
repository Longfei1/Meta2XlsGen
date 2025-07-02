package reader

import (
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
)

func WriteXls(path string, file *excelize.File) error {
	if err := os.MkdirAll(filepath.Dir(path), 644); err != nil {
		return err
	}

	return file.SaveAs(path)
}

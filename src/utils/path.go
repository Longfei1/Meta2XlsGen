package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetExecutablePath() string {
	exePath, err := os.Executable()
	if err != nil {
		return getExecutableByCaller()
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(res, tmpDir) {
		return getExecutableByCaller()
	}
	return res
}

func getExecutableByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = filepath.Dir(filename)
	}
	return abPath
}

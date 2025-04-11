package global

import (
	"os"
	"path/filepath"
)

var (
	BasePath, _ = os.Executable()
)

func init() {
	BasePath, _ = os.Executable()
	BasePath = filepath.Dir(BasePath)
}

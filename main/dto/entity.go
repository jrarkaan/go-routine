package dto

import (
	"os"
	"path/filepath"
)

var TempPath = filepath.Join(os.Getenv("TEMP"), "chapter-A.59-pipeline-temp")

package v2

import "path/filepath"

type Filesystem struct {
	target IFilesystemV2
}

var FilesystemApp Filesystem

func getRootPath(dir string) string {
	rootPath, _ := filepath.Abs(".")

	return filepath.Clean(filepath.Join(rootPath, dir))
}

// New 实例化
func (*Filesystem) New(target IFilesystemV2) *Filesystem { return &Filesystem{target: target} }

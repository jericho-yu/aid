package v2

import (
	"path/filepath"

	"github.com/jericho-yu/aid/array"
)

func getRootPath(dir string) string {
	rootPath, _ := filepath.Abs(".")

	return filepath.Clean(filepath.Join(rootPath, dir))
}

func CopyFiles(srcFiles, dstFiles *array.AnyArray[*File]) {
	srcFiles.Each(func(idx int, item *File) { item.CopyTo(dstFiles.Get(idx).GetFullPath()) })
}

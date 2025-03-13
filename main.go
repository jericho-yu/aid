package main

import (
	filesystem "github.com/jericho-yu/aid/filesystem/v2"
)

func main() {
	src := filesystem.DirApp.NewByRel("./a1")
	src.Ls()
	// dst := src.GetFiles().Copy().Every(func(item *filesystem.File) *filesystem.File {
	// 	return filesystem.FileApp.NewByAbs(filepath.Join(item.GetBasePath(), "..", "a2", item.GetName()))
	// })
	// dst.Each(func(idx int, item *filesystem.File) { fmt.Printf("%s\n", item.GetFullPath()) })

	filesystem.CopyFilesByDstPath(src.GetFiles(), src.Join("../a2").GetFullPath())
}

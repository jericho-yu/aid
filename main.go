package main

import (
	"log"

	filesystem "github.com/jericho-yu/aid/filesystem/v2"
	"github.com/jericho-yu/aid/httpClient"
)

func f1() {
	fs := filesystem.FileApp.NewByRel("./abc.txt")

	client := httpClient.App.NewGet("http://127.0.0.1:8080/download2").Download(fs.GetFullPath()).SaveLocal()
	if client.Err != nil {
		log.Fatalf("save local: %v", client.Err)
	}

	fsCtx := fs.Read()
	if fs.Error() != nil {
		log.Fatalf("read file: %v", fs.Error())
	}

	log.Printf("read file: %s\n", fsCtx)
}

func f2() {
	client := httpClient.App.NewGet("http://127.0.0.1:8080/ping").Send()
	log.Printf("ok: %s\n", client.GetResponseRawBody())
}

func main() {
	f1()
}

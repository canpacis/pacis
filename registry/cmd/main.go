package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed public/*
var public embed.FS

func main() {
	fs, _ := fs.Sub(public, "public")
	http.ListenAndServe(":8080", http.FileServerFS(fs))
}

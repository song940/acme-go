package main

import (
	"github.com/song940/acme-go/web"
)

func main() {
	server := web.NewServer()
	server.Start()
}

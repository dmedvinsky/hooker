package main

import (
	"fmt"
	"github.com/zenazn/goji"
)

func main() {
	LoadConfig()
	ConnectDB()
	goji.Get("/", Root)
	goji.Get("/hook/:guid", GetData)
	goji.Post("/hook/:guid", PostData)
	goji.Serve()
}

func hookUrl(guid string) string {
	return fmt.Sprintf("/hook/%s", guid)
}

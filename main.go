package main

import (
	"thor/src/server"
	"thor/src/server/container"
)

func main() {
	server.StartHttpServer(container.IntializeContainer())
}

package main

import (
	"github.com/Monstergogo/beauty-share/cmd/server/interface/rest"
	"github.com/Monstergogo/beauty-share/init/server"
)

func main() {
	s := server.InitServer()
	rest.InitRouter(s)
	s.RunServer()
}

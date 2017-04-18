package main

import (
	diamond "github.com/aerth/diamond/lib"
	"github.com/aertx/tlde/srv"
	"os"
)

var config = diamond.ConfigFields{
	Addr:     "0.0.0.0:" + os.Getenv("PORT"),
	Name:     "tl;de",
	Level:    3,
	Socket:   "/tmp/tlde.socket",
	Kicks:    true,
	Kickable: true,
}

func main() {
	handler := srv.NewMux()
	server := diamond.NewServer()
	server.Config = config
	server.SetMux(handler)
	err := server.Start()
	if err != nil {
		println(err.Error())
		os.Exit(111)
	}
	<-server.Done

}

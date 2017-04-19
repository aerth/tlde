package main

import (
	"os"

	diamond "github.com/aerth/diamond/lib"
	"github.com/aerth/tlde/src/tilde"
)

var config = diamond.ConfigFields{
	SocketHTTP:  os.Getenv("SOCKET"),
	Name:       "tl;de",
	Level:      3,
	Socket:     os.Getenv("ADMIN"),
	Kicks:      true,
	Kickable:   true,
}

func init(){
	println(tilde.Version())
	if os.Getenv("ADMIN") == "" ||
	   os.Getenv("PORT") == "" {
		println("need $ADMIN location and $PORT number")
		println("optional $SOCKET http unix socket location")
		println("example: env PORT=8080 ADMIN=./tlde.socket tlde")
	   	os.Exit(111)
	   }

	config.Addr = "0.0.0.0:" + os.Getenv("PORT")
}

func main() {
	server := diamond.NewServer(tilde.Handler())
	server.Config = config
	if os.Getenv("PORT") == "" {
		server.Config.Addr = ""
	}
	err := server.Start()
	if err != nil {
		println(err.Error())
		os.Exit(111)
	}
	<-server.Done

}

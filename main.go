// tlde web server
package main

import (
	"os"

	diamond "github.com/aerth/diamond/lib"
	"github.com/aerth/tlde/src/tilde"
)

var config = diamond.ConfigFields{
	SocketHTTP: os.Getenv("SOCKET"),
	Name:       "tl;de",
	Level:      3,
	Socket:     os.Getenv("ADMIN"),
	Kicks:      true,
	Kickable:   true,
}

func init() {
	println(tilde.Version())
	if os.Getenv("ADMIN") == "" {
		println("fatal: need location of admin socket. set $ADMIN")
		usage()
		os.Exit(111)
	}
	if len(os.Args) > 1 {
		usage()
		os.Exit(111)
	}
}

func usage() {
	println("USAGE")
	println("tlde listens on a unix socket or network port or both")
	println("specify listener by setting environmental variables PORT and/or SOCKET")
	println("specify location of admin socket with variable ADMIN")
	println()
	println("VARIABLES")
	println("Use environmental variables to modify how tlde works")
	println("\tADMIN\tlocation of diamond admin socket")
	println("\tADMIN\tlocation of diamond admin socket")
	println("\tSOCKET\tlocation of http listener socket")
	println("\tPORT\tport to listen on (see INTERFACE)")
	println("\tPUBLICDIR\tdirectory to serve in each user dir (default: 'Public')")
	println("\tFORMATPATH\tprintf style formatting for finding user dir (default: '/home/%s/%s')")
	println("\tINTERFACE\tinterface to listen on (default: 0.0.0.0)")

	println("\nEXAMPLE")
	println("env PORT=8080 ADMIN=./tlde.socket tlde")
	os.Exit(111)
}

func init() {
	if os.Getenv("PORT") == "" && os.Getenv("SOCKET") == "" {
		println("fatal: need either $PORT or $SOCKET (or both) to listen on")
		usage()
		os.Exit(111)
	}
	addr := os.Getenv("ADDRESS")
	if addr == "" {
		addr = "0.0.0.0"
	}
	config.Addr = addr + ":" + os.Getenv("PORT")
}

func main() {
	server := diamond.NewServer(tilde.NewHandler())
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

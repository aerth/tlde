// tlde web server
package main

import (
	"os"

	diamond "github.com/aerth/diamond/lib"
	"github.com/aerth/tlde/src/tilde"
)

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

func main() {
	if os.Getenv("PORT") == "" && os.Getenv("SOCKET") == "" {
		println("fatal: need either $PORT or $SOCKET (or both) to listen on")
		usage()
		os.Exit(111)
	}
	addr := os.Getenv("ADDRESS")
	if addr == "" {
		addr = "0.0.0.0"
	}
	if os.Getenv("PORT") != "" {
		addr = addr + ":" + os.Getenv("PORT")
	}

	server, err := diamond.NewServer(tilde.NewHandler(), os.Getenv("ADMIN"))
	if err != nil {
		println(err.Error())
		os.Exit(111)
	}
	server.Config.Kickable = true
	server.AddListener("tcp", addr)
	os.Clearenv()
	server.Runlevel(1)
	err = server.Runlevel(3)
	if err != nil {
		println(err.Error())
		os.Exit(111)
	}
	os.Exit(server.Wait())

}

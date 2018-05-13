package main

import (
	"github.com/teejays/copperchain"
	go_up "github.com/ufoscout/go-up"
	"log"
)

func main() {

	// Read environment variables with the help of go-up package
	up, err := go_up.NewGoUp().AddFile(".env", false).Build()
	if err != nil {
		log.Fatal(err)
	}

	dbRoot := up.GetString("db.root")
	serverHost := up.GetString("server.host")
	serverPort := up.GetInt("server.port")

	// Init Copper Chain
	copperchain.Init(copperchain.Options{
		DataRoot: dbRoot,
	})

	// (optional) Run Server
	err = copperchain.RunServer(copperchain.ServerOptions{
		Host: serverHost,
		Port: serverPort,
	})
	if err != nil {
		log.Fatal(err)
	}

}

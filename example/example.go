package main

import (
	"fmt"
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

	// Get CopperChain
	chain := copperchain.GetCopperChain()
	fmt.Printf("Chain: %+v\n", chain)

	// Add block data to chain (data should always be a map[string]interface{})
	var data map[string]interface{} = make(map[string]interface{})
	data["some_key"] = "the_value"
	err = copperchain.AddToCopperChain(data)
	if err != nil {
		log.Fatal(err)
	}
	// Check that the changes are reflected
	chain = copperchain.GetCopperChain()
	fmt.Printf("Chain: %+v\n", chain)

	// (optional) Run Server, and call http endpoints
	err = copperchain.RunServer(copperchain.ServerOptions{
		Host: serverHost,
		Port: serverPort,
	})
	if err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/teejays/copperchain"
	"github.com/ufoscout/go-up"
	"log"
)

var up go_up.GoUp

func main() {

	// Read environment variables with the help of go-up package
	var err error
	up, err = go_up.NewGoUp().AddFile(".env", false).Build()
	if err != nil {
		log.Fatal(err)
	}

	dbRoot := up.GetString("db.root")

	// Init Copper Chain
	copperchain.InitCopperChain(copperchain.Options{
		DataRoot: dbRoot,
	})

	// print some stuff
	dumpMyChain()

	// (optional) Run Tcp Server, and call http endpoints
	startExampleTcpServer()

}

func dumpMyChain() {
	// Get CopperChain
	chain := copperchain.GetMyChain()
	spew.Dump(chain)
}

func addDataToMyChain() {
	// Add block data to chain (data should always be a map[string]interface{})
	var data map[string]interface{} = make(map[string]interface{})
	data["some_key"] = "the_value"
	err := copperchain.AddToMyChain(data)
	if err != nil {
		log.Fatal(err)
	}
	// Check that the changes are reflected
	chain := copperchain.GetMyChain()
	spew.Dump(chain)
}

func startExampleTcpServer() {
	serverHost := up.GetString("tcp.server.host")
	serverPort := up.GetInt("tcp.server.port")
	err := copperchain.RunTcpServer(copperchain.ServerOptions{
		Host: serverHost,
		Port: serverPort,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func startExampleHttpServer() {
	serverHost := up.GetString("http.server.host")
	serverPort := up.GetInt("http.server.port")
	err := copperchain.RunHttpServer(copperchain.ServerOptions{
		Host: serverHost,
		Port: serverPort,
	})
	if err != nil {
		log.Fatal(err)
	}
}

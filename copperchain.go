package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/teejays/gofiledb"
	go_up "github.com/ufoscout/go-up"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var copperChain *BlockChain
var up go_up.GoUp

func main() {
	// Load environment variables
	var err error
	up, err = go_up.NewGoUp().AddFile("./.env", false).Build()
	if err != nil {
		log.Fatal(err)
	}

	// Set up DB client to store the block chain when the server is off
	dbRoot := up.GetStringOrDefault("db.root", ".data")
	gofiledb.InitClient(dbRoot)

	// Load the blockchain
	copperChain, err = LoadBlockChain()
	if err != nil {
		log.Fatal(err)
	}

	err = runServer()
	if err != nil {
		log.Fatal(err)
	}
}

func runServer() error {
	// Start the webserver
	port := up.GetIntOrDefault("server.port", 8080)
	addr := up.GetString("server.address")
	router := mux.NewRouter()
	router.HandleFunc("/", handleGetBlockChain).Methods("GET")
	router.HandleFunc("/", handleWriteBlock).Methods("POST")

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", addr, port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Printf("Server listening on %s:%d...\n", addr, port)

	return server.ListenAndServe()
}

// handleGetBlockChain listens for GET requests and serves the existing
// blockchain.
func handleGetBlockChain(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(copperChain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

// handleWriteBlock listens for POST requests with payload that is
// a valid BlockData. It creates a new Block for that data and adds
// it to the blockchain.
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var blockData BlockData
	err = json.Unmarshal(body, &blockData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the parent
	parent, err := copperChain.GetLastBlock(true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	block, err := NewBlock(blockData, parent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = copperChain.AddBlock(*block)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

package copperchain

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

type ServerOptions struct {
	Host string
	Port int
}

var defaultServerOptions ServerOptions = ServerOptions{
	Port: 8080,
}

// RunServer starts a webserver on the address provided and listens on two endpoints:
// 1. GET /: serves the block chain
// 2. POST /: adds data to the block chain. It accepts a json encoded BlockData as http body.
func RunServer(options ServerOptions) error {

	// Validate the parameters
	if options.Host == "" {
		fmt.Printf("empty host passed for server, Go will default to localhost.")
	}
	if options.Port < 1 {
		fmt.Printf("invalid port  '%d' passed in server, defaulting to port %d.", defaultServerOptions.Port)
		options.Port = defaultServerOptions.Port
	}

	// Start the webserver
	router := mux.NewRouter()
	router.HandleFunc("/", HandleGetBlockChain).Methods("GET")
	router.HandleFunc("/", HandleWriteBlock).Methods("POST")

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", options.Host, options.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Printf("Server listening on %s...\n", server.Addr)

	return server.ListenAndServe()
}

// HandleGetBlockChain listens for GET requests and serves the existing
// blockchain.
func HandleGetBlockChain(w http.ResponseWriter, r *http.Request) {
	// Get the BlockChain, json encode it and send it.
	chain := GetCopperChain()

	response, err := json.Marshal(chain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

// handleWriteBlock listens for POST requests with payload that is
// a valid BlockData. It creates a new Block for that data and adds
// it to the blockchain.
func HandleWriteBlock(w http.ResponseWriter, r *http.Request) {
	// parse the request body to get the blockdata
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

	// add the data to the block chain
	err = copperChain.AddBlockData(blockData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

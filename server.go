package copperchain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type ServerOptions struct {
	Host string
	Port int
}

const DEFAULT_HTTP_PORT int = 8080
const DEFAULT_TCP_PORT int = 9000

var defaultServerOptions ServerOptions = ServerOptions{
	Port: 8080,
}

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * *
* T C P  S E R V E R 			  						 *
* * * * * * * * * * * * * * * * * * * * * * * * * * * * */

var bcServer chan BlockChain

func RunTcpServer(options ServerOptions) error {
	bcServer = make(chan BlockChain)

	// Validate the parameters
	if options.Host == "" {
		fmt.Printf("empty host passed for server, Go will default to localhost.")
	}
	if options.Port < 1 {
		fmt.Printf("invalid port  '%d' passed in server, defaulting to port %d.", defaultServerOptions.Port)
		options.Port = DEFAULT_TCP_PORT
	}

	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", options.Host, options.Port))
	if err != nil {
		return err
	}
	defer server.Close()

	for {
		// everytime the server gets a new connection request
		conn, err := server.Accept()
		if err != nil {
			return err
		}
		go handleTcpConn(conn)
	}
}

func handleTcpConn(conn net.Conn) {
	defer conn.Close()

	io.WriteString(conn, "Please provide your data:")

	scanner := bufio.NewScanner(conn)

	// start up the process that sends myChain to connections every x seconds
	go func() {
		for {
			time.Sleep(30 * time.Second)
			myChainBytes, err := json.Marshal(myChain)
			if err != nil {
				log.Fatal(err)
			}
			io.WriteString(conn, string(myChainBytes))
		}
	}()

	// take in data from the TCP connection
	func() {
		for scanner.Scan() {
			fmt.Printf("Data received from %s", conn.LocalAddr())
			data := scanner.Bytes()
			var blockData BlockData
			err := json.Unmarshal(data, &blockData)
			if err != nil {
				log.Printf("Error while unmarshalling scanned data: %v", err)
				continue
			}
			err = AddToMyChain(blockData)
			if err != nil {
				log.Printf("Error while adding data to chain: %v", err)
				continue
			}

			bcServer <- myChain.Chain

			io.WriteString(conn, "Please provide your data:")
		}

	}()

}

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * *
* H T T P  S E R V E R 			  						 *
* * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// RunServer starts a webserver on the address provided and listens on two endpoints:
// 1. GET /: serves the blockchain
// 2. POST /: adds data to the blockchain. It accepts a json encoded BlockData as http body.
func RunHttpServer(options ServerOptions) error {

	// Validate the parameters
	if options.Host == "" {
		fmt.Printf("empty host passed for server, Go will default to localhost.")
	}
	if options.Port < 1 {
		fmt.Printf("invalid port  '%d' passed in server, defaulting to port %d.", defaultServerOptions.Port)
		options.Port = DEFAULT_HTTP_PORT
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

// HandleGetBlockChain listens for GET requests and response with the
// blockchain. This is basically a HTTP getter for blockchain.
func HandleGetBlockChain(w http.ResponseWriter, r *http.Request) {
	// Get the BlockChain, json encode it and send it.
	chain := GetMyChain()

	response, err := json.Marshal(chain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

// HandleWriteBlock listens for POST requests with payload that is
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

	// add the data to the blockchain
	err = AddToMyChain(blockData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

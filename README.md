# CopperChain

CopperChain implements a basic blockchain in Golang. It treats blockchain as a package, which can be imported by other packages and used. It also has the optional functionality of a web server, which provides basic Get and Post endpoints.

Click [here](https://godoc.org/github.com/teejays/copperchain) for code documentation.
 
## Getting Started

### Prerequisites:
1) Install Golang. You can install it from the [official Go website](https://golang.org/).
2) Install the following Go packages (if you don't have them already):
	* [GoFiledb](https://github.com/teejays/gofiledb): A micro, very easy to use, DB client that uses filesystem for storage. To install, run: `go get -u github.com/teejays/gofiledb`
	* [Gorilla Mux](https://github.com/gorilla/mux): A powerful URL router and dispatcher for golang. To install, run: ` go get -u github.com/gorilla/mux`
    
    
### Usage:

First, install CopperChain by running in your command line:

`go get -u github.com/teejays/copperchain`

After this, CopperChain can be included in your package easily. Add the following import statement to the go file which needs to use the package:
```go
import github.com/teejays/copperchain
```
Before you interact with CopperChain, you need to initialize it. You can initialize it by calling the Init() function, that accepts optional parameters.

```go
// Init Copper Chain
copperchain.Init(copperchain.Options{})
```

Once initialized, you can get the existing blockchain by calling the GetCopperChain() method.
```go
chain := copperchain.GetCopperChain()
```

You can add a block to the blockchain by calling the AddToCopperChain() method, which requires a parameter of the form map[string]interface. This allows for flexibility of data that can be kept in the block.
```go
data := make(map[string]interface{})
data["transaction_id"] = "abc123"
err := copperchain.AddToCopperChain(data)
```

If you want to set up the blockchain as RESTful API, you can call the RunSerever() method, which takes some options around the server address on which to rum the server. 
```go
err := copperchain.RunServer(copperchain.ServerOptions{
		Host: "127.0.0.1", //optional
		Port: 8080, // default 8080
	})
```
#### Example Code
Here is an example code that demonstrates the use of CopperChain. It can also be found at [example.go](https://github.com/teejays/copperchain/example).
```go
package main

import (
	"fmt"
	"github.com/teejays/copperchain"
	"log"
)

func main() {

	// Init CopperChain
	copperchain.Init(copperchain.Options{
		DataRoot: ".data", 
	})
	
    	// Get CopperChain
	chain := copperchain.GetCopperChain()
	fmt.Printf("Chain: %+v\n", chain)

	// Add block data to chain (data should always be a map[string]interface{})
	var data map[string]interface{} = make(map[string]interface{})
	data["some_key"] = "the_value"
	copperchain.AddToCopperChain(data)

	// Check that the changes are reflected
	chain = copperchain.GetCopperChain()
	fmt.Printf("Chain: %+v\n", chain)
    
	// (optional) Run Server
	err = copperchain.RunServer(copperchain.ServerOptions{
		Host: serverHost,
		Port: serverPort,
	})
	if err != nil {
		log.Fatal(err)
	}

}
```

To manually run the provided example, download the example code locally, build it and run it:
```
go build -o example
./example
```

#### API Endpoints
If you decide to run the built-in webserver, it will listen on two endpoints:
1. **HTTP GET /** 
	* provides the entire chain
	* _curl example_: `curl http://localhost:8080/`  
2. **HTTP POST /**
	* accepts a json object as body, which can be parsed into a `map[string]interface{}`
	* _curl example_: `curl http://localhost:8080/ -H "Content-Type: application/json" -d '{"transaction_id":"abc123"}'`
	
  ### Contact
  
  For any issues, please open a new issue. If you want to contribute, please feel free to submit a merge request or reach out to me at copperchain@teejay.me.

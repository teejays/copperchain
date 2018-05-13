// Copperchain implements a basic blockchain in GoLang. It provides the
// blockchain as a package, with an optional functionlaity to starting a server.
// Sample usage of copperchain package is provided in github.com/teejays/copperchain/example
package copperchain // import "github.com/teejays/copperchain"

import (
	"fmt"
	"github.com/teejays/gofiledb"
	"log"
)

// copperChain is the primary blockchain that we keep in memory as the program runs.
var copperChain *BlockChainAtomic

// Options represents that struct of the parameter that should be passed
// when initialing the CopperChain.
type Options struct {
	DataRoot string
}

// defaultOptions specify the default parameters for Init().
var defaultOptions Options = Options{
	DataRoot: ".data",
}

// Init initializes the CopperChain package. It requires Options as a parameter.
// It does a few things: 1) it initializes the GoFiledb for local storage of the
// bloch chain, 2) it intializes the copperChain global variable, populating it
// with previously saved copperChain data in the GoFiledb, if any. It panics if it
// encounters an error.
func Init(options Options) {

	// Validate the options, and resort to default when needed
	if options.DataRoot == "" {
		fmt.Printf("empty DataRoot passed in options for copperchain, defaulting to %s.", defaultOptions.DataRoot)
		options.DataRoot = defaultOptions.DataRoot
	}

	// Initiate GoFiledb so blockchain instances can be saved
	gofiledb.InitClient(options.DataRoot)

	// Read the saved blockchain using GoFileDb and put as the global chain
	var newCopperChain BlockChainAtomic
	var err error
	newCopperChain.Chain, err = loadBlockChainFromDb()
	if err != nil {
		log.Panic(err)
	}
	copperChain = &newCopperChain

}

// GetCopperChain is a getter method to access the blockchain. This is the only exported
// method that allows access to the chain from outside this package.
func GetCopperChain() BlockChain {
	return copperChain.Chain
}

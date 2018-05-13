package copperchain // import "github.com/teejays/copperchain"

import (
	"fmt"
	"github.com/teejays/gofiledb"
	"log"
)

var copperChain *BlockChainAtomic

type Options struct {
	DataRoot string
}

var defaultOptions Options = Options{
	DataRoot: ".data",
}

var isInitiated bool

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

	// set this variable so other functions can know that the package has been initiated properly
	isInitiated = true

}

func GetCopperChain() BlockChain {
	return copperChain.Chain
}

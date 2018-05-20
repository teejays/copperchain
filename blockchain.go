package copperchain

import (
	"fmt"
	"github.com/teejays/gofiledb"
	"sync"
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * *
* B L O C K   C H A I N  								 *
* * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// BlockChain is made of a slice of Blocks.
type BlockChain []Block

// BlockChainAtomic type implements a thread safe version of a BlockChain. All the
// struct methods for BlockChain are implemented on BlockChainAtomic.
type BlockChainAtomic struct {
	Chain BlockChain
	sync.RWMutex
}

func (chain *BlockChainAtomic) Length() int {
	return len(chain.Chain)
}

// // AddBlock adds a new block into the blockchain with the provided blockData. In the process is
// // runs some validation to ensure the integrity of the blockchain.
// func (chain *BlockChainAtomic) AddBlockData(blockData BlockData) error {

// 	// get the parent first block first
// 	parent, err := copperChain.GetLastBlock(true)
// 	if err != nil {
// 		return err
// 	}

// 	// create the new Block
// 	block, err := newBlock(blockData, parent)
// 	if err != nil {
// 		return err
// 	}

// 	chain.AppendBlock(block)

// 	return nil
// }

func (chain *BlockChainAtomic) AppendBlock(block Block) error {
	chain.Lock()
	defer chain.Unlock()
	newChain := append(chain.Chain, block)

	err := block.validateBlockWithParent(newChain[len(newChain)-2])
	if err != nil {
		return err
	}
	chain.Chain = newChain

	return nil
}

// IsBlockValid takes the index representing the location of a block in
// the chain, and checks whether that block has valid fields, and
// adheres to it's relationship with its parent block.
func (chain *BlockChainAtomic) ValidateBlockAtIndex(index int) error {

	block, err := chain.GetBlockByIndex(index, true)
	if err != nil {
		return err
	}

	err = block.validateFields()
	if err != nil {
		return err
	}

	previousBlock, err := chain.GetBlockByIndex(index-1, true)
	if err != nil {
		return err
	}

	err = block.validateBlockWithParent(previousBlock)
	if err != nil {
		return err
	}

	return nil
}

// GetLastBlock provides the last block in the blockchain. This is
// usually the parent of any incoming new block. Use useLock as true in
// order to execute this function in a threadsafe way.
func (chain *BlockChainAtomic) GetLastBlock(useLock bool) (Block, error) {
	var block Block
	if chain.Length() == 0 {
		return block, fmt.Errorf("tThere are no blocks in the chain")
	}
	return chain.GetBlockByIndex(chain.Length()-1, useLock)
}

// GetBlockByIndex provides the block tat resides at the given index in
// the blockchain. Use useLock as true in  order to execute this
// function in a threadsafe way.
func (chain *BlockChainAtomic) GetBlockByIndex(index int, useLock bool) (Block, error) {
	var block Block
	if index < 0 {
		return block, fmt.Errorf("index provided for GetBlockByIndex '%d' is not valid", index)
	}
	if index >= len(chain.Chain) {
		return block, fmt.Errorf("index provided for GetBlockByIndex '%d' is greater then the length of the blockchain '%d'", index, len(chain.Chain))
	}
	if useLock {
		chain.RLock()
	}
	block = chain.Chain[index]
	if useLock {
		chain.RUnlock()
	}
	return block, nil
}

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * *
* H E L P E R S 		 								 *
* * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// loadBlockChainFromDb reads the already saved blockchain data from the
// database. We use this function to make sure that we can use the
// previously saved instance of the chain upon startup.
func loadBlockChainFromDb() (BlockChain, error) {
	var chain BlockChain = make([]Block, 0) // initialize to an empty array vs. null since null sounds wierd
	db := gofiledb.GetClient()
	_, err := db.GetStructIfExists("blockchain", "blockchain_v1", &chain)
	if err != nil {
		return nil, err
	}
	return chain, nil
}

// saveBlockChainToDb saves the instance of blockchain into the database.
func saveBlockChainToDb(chain BlockChain) error {
	db := gofiledb.GetClient()
	return db.SetStruct("blockchain", "blockchain_v1", chain)
}

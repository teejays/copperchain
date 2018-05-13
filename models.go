package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/teejays/gofiledb"
	"strconv"
	"sync"
	"time"
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * *
* B L O C K   C H A I N  								 *
* * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// BlockChain type implements a thread safe version of a chain.
// The blocks are stored in the Chain field.
type BlockChain struct {
	Chain []Block
	Lock  sync.RWMutex
}

// LoadBlockChain reads the already saved blockchain data from the
// database. We use this function to make sure that we can use the
// previously saved instance of the chain upon startup.
func LoadBlockChain() (*BlockChain, error) {
	var chain BlockChain
	db := gofiledb.GetClient()
	_, err := db.GetStructIfExists("blockchain", "blockchain_v1", &chain)
	if err != nil {
		return nil, err
	}
	return &chain, nil
}

// Save the instance of blockchain into the database.
func (chain *BlockChain) Save() error {
	db := gofiledb.GetClient()
	return db.SetStruct("blockchain", "blockchain_v1", chain)
}

// GetLastBlock provides the last block in the block chain. This is
// usually the parent of any incoming new block. Use useLock as true in
// order to execute this function in a threadsafe way.
func (chain *BlockChain) GetLastBlock(useLock bool) (*Block, error) {
	lenChain := len(chain.Chain)
	if lenChain == 0 {
		return nil, nil
	}
	return chain.GetBlockByIndex(lenChain-1, useLock)
}

// GetBlockByIndex provides the block tat resides at the given index in
// the blockchain. Use useLock as true in  order to execute this
// function in a threadsafe way.
func (chain *BlockChain) GetBlockByIndex(index int, useLock bool) (*Block, error) {
	if index < 0 {
		return nil, fmt.Errorf("index provided for GetBlockByIndex '%d' is not valid", index)
	}
	if index >= len(chain.Chain) {
		return nil, fmt.Errorf("index provided for GetBlockByIndex '%d' is greater then the length of the block chain '%d'", index, len(chain.Chain))
	}
	if useLock {
		chain.Lock.RLock()
	}
	block := chain.Chain[index]
	if useLock {
		chain.Lock.RUnlock()
	}
	return &block, nil
}

// AddBlock adds a new block into the blockchain. In the process is
// runs some validation to ensure the integrity of the blockchain.
func (chain *BlockChain) AddBlock(block Block) error {

	// validate the block fields
	err := block.ValidateFields()
	if err != nil {
		return err
	}

	chain.Lock.Lock()
	defer chain.Lock.Unlock()

	parent, err := chain.GetLastBlock(false)
	if err != nil {
		return err
	}

	err = block.ValidateBlockWithParent(parent)
	if err != nil {
		return err
	}

	chain.Chain = append(chain.Chain, block)
	chain.Save()

	return nil
}

// IsBlockValid takes the index representing the location of a block in
// the chain, and checks whether that block has valid fields, and
// adheres to it's relationship with its parent block.
func (chain *BlockChain) ValidateBlockAtIndex(index int) error {

	block, err := chain.GetBlockByIndex(index, true)
	if err != nil {
		return err
	}

	err = block.ValidateFields()
	if err != nil {
		return err
	}

	parent, err := chain.GetBlockByIndex(index-1, true)
	if err != nil {
		return err
	}

	err = block.ValidateBlockWithParent(parent)
	if err != nil {
		return err
	}

	return nil
}

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * *
* B L O C K 			  								 *
* * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// Block is the primary unit of a BlockChain.
type Block struct {
	Index     int
	Timestamp time.Time
	Data      BlockData
	Hash      string
	PrevHash  string
}

// BlockData represents the structure in which data in the blockchain can
// stored. It is defined as a map[string]interface{} to allow for
// flexibility on what and how much data could be contained in the block.
type BlockData map[string]interface{}

// NewBlock returns an instance of a Block initialized with the provided
// data and the parent block.
func NewBlock(data BlockData, parent *Block) (*Block, error) {
	// verify that the data is valid for the new block
	if data == nil {
		return nil, fmt.Errorf("attempted to create a new block with nil data")
	}

	// create a new block
	var b Block
	b.Timestamp = time.Now()
	b.Data = data
	if parent != nil {
		b.Index = parent.Index + 1
		b.PrevHash = parent.Hash
	}

	b.Hash = b.calculateHash()

	return &b, nil
}

// calculateHash takes into account all the fields of the block to
// create a string representation of a hash that is unique to it.
func (b *Block) calculateHash() string {
	// create a string representation of the block
	// one way to do that is to concatenate the string representation of the individual fields
	string_record := strconv.Itoa(b.Index) + strconv.FormatInt(b.Timestamp.UnixNano(), 10) + fmt.Sprintf("%#v", b.Data) + b.PrevHash

	// create a hash
	h := sha256.New()
	h.Write([]byte(string_record))
	hashed := h.Sum(nil)
	return string(hashed)
}

// ValidateBlockWithParent checks whether a block adheres to it's
// relationship with its parent block.
func (b *Block) ValidateBlockWithParent(parent *Block) error {
	if parent == nil {
		fmt.Println("Request made to ValidateBlockWithParent with a nil parent. Is it the first ever block in the chain?")
		return nil
	}
	if b.Index != parent.Index+1 {
		return fmt.Errorf("index '%d' is not equal to 1 + index of parent '%d'", b.Index, parent.Index)
	}
	if b.PrevHash != parent.Hash {
		return fmt.Errorf("hash does not match parent hash")
	}
	return nil
}

// ValidateFields ensures that the fields of a block have valid values.
func (b *Block) ValidateFields() error {
	var errors []error
	if b.Index < 0 {
		errors = append(errors, fmt.Errorf("invalid index field %d", b.Index))
	}
	if b.Timestamp.IsZero() {
		errors = append(errors, fmt.Errorf("empty Timestamp field"))
	}
	if b.Data == nil {
		errors = append(errors, fmt.Errorf("nil Data field"))
	}
	if b.Hash == "" {
		errors = append(errors, fmt.Errorf("empty hash field"))
	}
	if b.Hash != b.calculateHash() {
		errors = append(errors, fmt.Errorf("unexpected hash"))
	}

	if len(errors) < 1 {
		return nil
	}

	var err_message string = fmt.Sprintf("Block validation failed with %d errors:", len(errors))
	for i, err := range errors {
		err_message += fmt.Sprintf(" (%d) %v", i, err)
	}
	return fmt.Errorf(err_message)
}

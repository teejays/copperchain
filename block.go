package copperchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * *
* B L O C K 			  								 *
* * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// BlockData represents the structure in which data in the blockchain can
// stored. It is defined as a map[string]interface{} to allow for
// flexibility on what and how much data could be contained in the block.
type BlockData map[string]interface{}

// Block is the primary unit of a BlockChain. It is linked to the previous block
// in the blockchain using the PrevHash field, which is the Hash field of its parent
// block.
type Block struct {
	Index     int
	Timestamp time.Time
	Data      BlockData
	Hash      string
	PrevHash  string
}

// newBlock returns an instance of a Block initialized with the provided
// data and the parent block.
func newBlock(data BlockData, parent *Block) (*Block, error) {
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
	return hex.EncodeToString(hashed)
}

// validateBlockWithParent checks whether a block adheres to it's
// relationship with its parent block.
func (b *Block) validateBlockWithParent(parent *Block) error {
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

// validateFields ensures that the fields of a block have valid values.
func (b *Block) validateFields() error {
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

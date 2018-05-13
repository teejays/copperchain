package copperchain

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type BlockChain []Block

// IsBlockValid takes the index representing the location of a block in the chain, and checks whether that block is valid.
func (chain BlockChain) IsBlockValid(index int) bool {
	block := chain[index]
	if block.Index != index {
		return false
	}
	// if it's the first element
	if index < 1 {
		return true
	}
	// else, we should chech block's integrity respecctive to it's parent
	parent := chain[index-1]
	if block.PrevHash != parent.Hash {
		return false
	}
	if block.calculateHash() != block.Hash {
		return false
	}
	return true
}

type Block struct {
	Index     int
	Timestamp time.Time
	Data      interface{}
	Hash      string
	PrevHash  string
}

func NewBlock(parent Block, data interface{}) (*Block, error) {
	// verify that parent block is valid
	err := parent.ValidateFields()
	if err != nil {
		return nil, err
	}
	// verify that the data is valid for the new block
	if data == nil {
		return nil, fmt.Errorf("attempted to create a new block with nil data")
	}

	// create a new block
	var b Block
	b.Index = parent.Index + 1
	b.Timestamp = time.Now()
	// do not store a pointer I guess, because pointer data can be changed outside the block chain
	b.Data = data
	b.PrevHash = parent.Hash
	b.Hash = b.calculateHash()

	return &b, nil
}

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
	// if b.PrevHash == "" {
	// 	errors = append(errors, fmt.Error("empty previous hash field"))
	// }

	if len(errors) < 1 {
		return nil
	}

	var err_message string = fmt.Sprintf("Block validation failed with %d errors:", len(errors))
	for i, err := range errors {
		err_message += fmt.Sprintf(" (%d) %v", i, err)
	}
	return fmt.Errorf(err_message)
}

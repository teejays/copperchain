package copperchain

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Index     int64
	Timestamp time.Time
	Data      interface{}
	Hash      string
	PrevHash  string
}

type BlockChain []Block

func main() {

	return
}

func (b *Block) calculateHash() string {
	// create a string representation of the block
	// one way to do that is to concatenate the string representation of the individual fields
	string_record := strconv.FormatInt(b.Index, 10) + strconv.FormatInt(b.Timestamp.UnixNano(), 10) + fmt.Sprintf("%#v", b.Data) + b.PrevHash

	// create a hash
	h := sha256.New()
	h.Write([]byte(string_record))
	hashed := h.Sum(nil)

	return string(hashed)
}

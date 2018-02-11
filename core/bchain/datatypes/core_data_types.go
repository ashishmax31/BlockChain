package datatypes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type Block struct {
	Index        int           `json:"index" `
	Hash         string        `json:"hash"`
	Timestamp    time.Time     `json:"created_at"`
	PrevHash     string        `json:"previous_hash"`
	Nounce       int64         `json:"nounce"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Signature string `json:"signature"`
	Amount    int64  `json:"amount"`
}

type BlockChain struct {
	sync.RWMutex
	BlkChain []Block
}

type CurrentTransactions struct {
	sync.RWMutex
	Transactions []Transaction
}

func (b *Block) SetHash() string {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%v", b)))
	sha256Hash := hex.EncodeToString(hasher.Sum(nil))
	b.Hash = sha256Hash
	return sha256Hash
}

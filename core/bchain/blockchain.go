package bchain

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
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   int64  `json:"amount"`
}

type blockChain struct {
	sync.RWMutex
	blkChain []Block
}

type currentTransactions struct {
	sync.RWMutex
	transactions []Transaction
}

var cTransactions currentTransactions

// BlockChain ... The actual in-memory immutable sequence of records called blocks.
var BlockChain blockChain

func init() {
	genesis()
}

func (b *blockChain) blockChainLen() int {
	b.RLock()
	defer b.RUnlock()
	return len(b.blkChain)
}

func (b *blockChain) ReadBlockChain() []Block {
	b.RLock()
	defer b.RUnlock()
	return b.blkChain
}

func NewTransaction(t Transaction) int {
	cTransactions.Lock()
	defer cTransactions.Unlock()
	cTransactions.transactions = append(cTransactions.transactions, t)
	return lastBlock().Index + 1
}

func lastBlock() Block {
	blkChainLen := BlockChain.blockChainLen()
	if blkChainLen > 0 {
		return BlockChain.blkChain[BlockChain.blockChainLen()-1]
	}
	// return a dummy block with hash `1` when the block chain is just born
	return Block{Hash: "1"}
}

func (b *blockChain) NewBlock(nounce int64) Block {
	newBlock := Block{
		Index:        b.blockChainLen() + 1,
		PrevHash:     lastBlock().Hash,
		Nounce:       nounce,
		Transactions: cTransactions.transactions,
		Timestamp:    time.Now(),
	}
	newBlock.setHash()
	b.Lock()
	defer b.Unlock()
	cTransactions.transactions = []Transaction{}
	b.blkChain = append(b.blkChain, newBlock)
	return newBlock
}

func (b *Block) setHash() string {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%v", b)))
	sha256Hash := hex.EncodeToString(hasher.Sum(nil))
	b.Hash = sha256Hash
	return sha256Hash
}

func genesis() Block {
	b := BlockChain.NewBlock(100)
	return b
}

func proofOfWork(prevProof int64) int64 {
	var proof int64
	for {
		if validProof(prevProof, proof) == true {
			break
		}
		proof++
	}
	// println(proof)
	return proof
}

func validProof(prevProof int64, proof int64) bool {
	str := fmt.Sprintf("%v%v", proof, prevProof)
	hasher := sha256.New()
	hasher.Write([]byte(str))
	sha1Hash := hex.EncodeToString(hasher.Sum(nil))
	// println(sha1Hash)
	if sha1Hash[:4] == "0000" {
		return true
	}
	return false
}

func Mine() {
	lastBlck := lastBlock()
	lastProof := lastBlck.Nounce
	proof := proofOfWork(lastProof)
	BlockChain.NewBlock(proof)
}

// func main() {
// 	genesis()
// 	t1 := transaction{Sender: "1", Receiver: "ashish", Amount: 10}
// 	t2 := transaction{Sender: "ashish", Receiver: "akshay", Amount: 5}
// 	newTransaction(t1)
// 	newTransaction(t2)
// 	fmt.Printf("%v\n", BlockChain.blkChain)
// 	// proofOfWork(0)
// 	mine()
// 	fmt.Printf("%v\n", BlockChain.blkChain)
// 	t2 = transaction{Sender: "ashish", Receiver: "akshay", Amount: 5}
// 	t1 = transaction{Sender: "1", Receiver: "ashish", Amount: 10}
// 	newTransaction(t1)
// 	newTransaction(t2)
// 	mine()
// 	fmt.Printf("%v\n", BlockChain.blkChain)

// }

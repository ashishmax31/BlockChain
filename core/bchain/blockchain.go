package bchain

import (
	"log"
	"time"

	"github.com/ashishmax31/blockchain/core/bchain/datatypes"
	"github.com/ashishmax31/blockchain/core/bchain/hashproblem"
	"github.com/ashishmax31/blockchain/core/consensus"
)

var cTransactions datatypes.CurrentTransactions

// BlockChain ... The actual in-memory immutable sequence of records called blocks.

type blockChain datatypes.BlockChain

var BlockChain blockChain

func init() {
	genesis()
}

func (b *blockChain) blockChainLen() int {
	b.RLock()
	defer b.RUnlock()
	return len(b.BlkChain)
}

func (b *blockChain) ReadBlockChain() []datatypes.Block {
	b.RLock()
	defer b.RUnlock()
	return b.BlkChain
}

func ReadBlockChain() []datatypes.Block {
	BlockChain.RLock()
	defer BlockChain.RUnlock()
	return BlockChain.BlkChain
}

func NewTransaction(t datatypes.Transaction) int {
	cTransactions.Lock()
	defer cTransactions.Unlock()
	cTransactions.Transactions = append(cTransactions.Transactions, t)
	return lastBlock().Index + 1
}

func lastBlock() datatypes.Block {
	blkChainLen := BlockChain.blockChainLen()
	if blkChainLen > 0 {
		return BlockChain.BlkChain[BlockChain.blockChainLen()-1]
	}
	// return a dummy block with hash `1` when the block chain is just born
	return datatypes.Block{Hash: "1"}
}

func (b *blockChain) NewBlock(nounce int64) datatypes.Block {
	newBlock := datatypes.Block{
		Index:        b.blockChainLen() + 1,
		PrevHash:     lastBlock().Hash,
		Nounce:       nounce,
		Transactions: cTransactions.Transactions,
		Timestamp:    time.Now(),
	}
	newBlock.SetHash()
	b.Lock()
	defer b.Unlock()
	cTransactions.Transactions = []datatypes.Transaction{}
	b.BlkChain = append(b.BlkChain, newBlock)
	return newBlock
}

func genesis() datatypes.Block {
	b := BlockChain.NewBlock(100)
	return b
}

func proofOfWork(prevProof int64) int64 {
	var proof int64
	for {
		if hashproblem.ValidProof(prevProof, proof) == true {
			break
		}
		proof++
	}
	// println(proof)
	return proof
}

func Mine() {
	cTransactions.RLock()
	defer cTransactions.RUnlock()
	if len(cTransactions.Transactions) == 0 {
		// Nothing to add to the block as no transactions have happened.
		// Not explicitly passing an error saying no transactions to mine, just silently returning.
		return
	}
	lastBlck := lastBlock()
	lastProof := lastBlck.Nounce
	proof := proofOfWork(lastProof)
	BlockChain.NewBlock(proof)
	resp, ind := consensus.ValidateChain(BlockChain.ReadBlockChain())
	if resp == true {
		log.Println("Block validated")
	} else {
		log.Println(ind)
	}
}

func ReplaceBlockChain(newBlock []datatypes.Block) {
	// First validate the incoming block chain, then update.
	resp, _ := consensus.ValidateChain(newBlock)
	if resp == true {
		// Secure with mutual exclusion, don't want to fuck up.
		BlockChain.Lock()
		defer BlockChain.Unlock()
		BlockChain.BlkChain = newBlock
	} else {
		log.Println("Recieved blockchain not valid!")
	}
}

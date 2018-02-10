package crons

import (
	"log"
	"time"

	"github.com/ashishmax31/blockchain/core/bchain"
	"github.com/ashishmax31/blockchain/core/consensus"

	"github.com/ashishmax31/blockchain/core/bchain/datatypes"
	"github.com/ashishmax31/blockchain/core/nodes"
)

func initiateConsensus() {
	log.Println("Timer worked!")
	if len(nodes.NodeList()) == 1 {
		return
	}
	log.Println("Initiating consensus protocol")
	nodesList := nodes.NodeList()
	longestestChain := struct {
		chain       []datatypes.Block
		nodeAddress string
	}{bchain.ReadBlockChain(), nodes.CurrentNodeAddress}
	for nodeAddress := range nodesList {
		if nodeAddress != nodes.CurrentNodeAddress {
			chain := consensus.GetChain(nodeAddress)
			if len(chain) > len(longestestChain.chain) {
				longestestChain.chain = chain
				longestestChain.nodeAddress = nodeAddress
			}
		}
	}
	// If the new chain is longest chain is recieved from another node, replace the current node's
	// chain.
	if longestestChain.nodeAddress != nodes.CurrentNodeAddress {
		log.Printf("Got new longest chain from: %s \n", longestestChain.nodeAddress)
		bchain.ReplaceBlockChain(longestestChain.chain)
	}
}

func pollOtherNodes() {
	timer := time.NewTicker(10 * time.Second)
	for {
		<-timer.C
		initiateConsensus()
	}
}

func init() {
	go pollOtherNodes()
}

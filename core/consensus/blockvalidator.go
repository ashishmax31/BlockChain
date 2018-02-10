package consensus

import (
	"github.com/ashishmax31/blockchain/core/bchain/datatypes"
	"github.com/ashishmax31/blockchain/core/bchain/hashproblem"
)

func validateChainHash(chain []datatypes.Block) (bool, int) {
	for index, block := range chain {
		if block.Index != len(chain) {
			if block.Hash != chain[index+1].PrevHash {
				return false, index
			}
		}
	}
	return true, -1
}

func validateProofOfWork(chain []datatypes.Block) (bool, int) {
	for index, block := range chain {
		if block.Index != 1 {
			if !hashproblem.ValidProof(chain[index-1].Nounce, block.Nounce) {
				return false, index
			}
		}
	}
	return true, -1
}

func ValidateChain(chain []datatypes.Block) (bool, int) {
	status, index := validateChainHash(chain)
	if status != true {
		return false, index
	}
	status, index = validateProofOfWork(chain)
	if status != true {
		return false, index
	}
	return true, -1
}

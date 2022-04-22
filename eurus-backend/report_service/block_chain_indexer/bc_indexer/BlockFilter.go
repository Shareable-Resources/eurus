package bc_indexer

import "github.com/ethereum/go-ethereum/core/types"

func BlockHasTransaction(block *types.Block)(bool){
	if(block.Transactions().Len()>0){
		return true
	}
	return false
}
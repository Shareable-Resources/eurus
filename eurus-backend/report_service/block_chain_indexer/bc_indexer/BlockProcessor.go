package bc_indexer

import (
	"github.com/ethereum/go-ethereum/core/types"
)

func ProcessBlock(context *blockChainProcessorContext, block *types.Block) {
	if BlockHasTransaction(block) {
		TransactionsHandler(context, block)
	}

}

func ProcessMainNetBlock(context *blockChainProcessorContext, block *types.Block) {
	if BlockHasTransaction(block) {
		TransactionsHandler(context, block)
	}

}

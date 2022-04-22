package ethereum

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockBuffer struct {
	Root     *BlockBufferNode
	MaxDepth int
	NodeMap  map[string]*BlockBufferNode //Block hash to node map
}

func NewBlockBuffer(maxDepth int) *BlockBuffer {
	blockBuffer := new(BlockBuffer)
	blockBuffer.NodeMap = make(map[string]*BlockBufferNode)
	blockBuffer.MaxDepth = maxDepth
	blockBuffer.Root = new(BlockBufferNode)

	return blockBuffer
}

func (me *BlockBuffer) AddNode(blockHeader *types.Header) bool {
	parent, ok := me.NodeMap[blockHeader.ParentHash.Hex()]
	child := NewBlockBufferNode()
	child.BlockHeader = blockHeader
	if !ok {
		if me.Root.BlockHeader == nil && (len(me.Root.Children) == 0 || me.Root.Children[0].BlockHeader.Number.Cmp(blockHeader.Number) == 0) {
			//Initial case
			me.Root.AddChild(child)
		} else {
			return false
		}
	} else {
		parent.AddChild(child)
	}
	me.NodeMap[child.Hash().Hex()] = child
	return true
}

func (me *BlockBuffer) GetConfirmedBlock() (int, *BlockBufferNode) {
	count, child := me.Root.GetLongestChild()
	if count < me.MaxDepth {
		return count, nil
	}
	me.TrimNodeMap(me.Root, child)
	me.Root.Reset()
	me.Root = child //The child becomes current confirmed node
	child.Parent = nil
	return count, child
}

func (me *BlockBuffer) TrimNodeMap(node *BlockBufferNode, excludeChild *BlockBufferNode) {

	for _, discardChild := range node.Children {
		if discardChild != excludeChild {
			delete(me.NodeMap, discardChild.Hash().Hex())
			me.TrimNodeMap(discardChild, excludeChild)
		}
	}
}

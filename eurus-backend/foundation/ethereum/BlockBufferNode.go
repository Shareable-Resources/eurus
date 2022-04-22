package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockBufferNode struct {
	Parent      *BlockBufferNode
	Children    []*BlockBufferNode
	BlockHeader *types.Header
}

func NewBlockBufferNode() *BlockBufferNode {
	node := new(BlockBufferNode)
	node.Children = make([]*BlockBufferNode, 0)
	return node
}

func (me *BlockBufferNode) AddChild(child *BlockBufferNode) {
	me.Children = append(me.Children, child)
	child.Parent = me
}

func (me *BlockBufferNode) GetLongestChild() (int, *BlockBufferNode) {
	var longestCount int = 1
	var longestNode *BlockBufferNode = nil
	for _, node := range me.Children {
		if node.Children != nil && len(node.Children) > 0 {
			count, _ := node.GetLongestChild()
			if count+1 > longestCount {
				longestCount = count + 1
				longestNode = node
			}
		}
	}

	return longestCount, longestNode
}

func (me *BlockBufferNode) Reset() {
	me.BlockHeader = nil
	me.Parent = nil
	me.Children = nil
}

func (me *BlockBufferNode) Hash() common.Hash {
	return me.BlockHeader.Hash()
}

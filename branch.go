package mptrie

import "github.com/ethereum/go-ethereum/crypto"

type BranchNode struct {
	Branches [16]Node
	Value    []byte
}

func NewBranchNode() *BranchNode {
	return &BranchNode{
		Branches: [16]Node{},
	}
}

func (b *BranchNode) SetBranch(nb Nibble, n Node) {
	b.Branches[int(nb)] = n
}

func (b *BranchNode) SetValue(v []byte) {
	b.Value = v
}

func (b BranchNode) Raw() []interface{} {
	hashes := make([]interface{}, 17)
	for i := 0; i < 16; i++ {
		if b.Branches[i] == nil {
			hashes[i] = EmptyNodeRaw
		} else {
			node := b.Branches[i]
			if len(Serialize(node)) >= 32 {
				hashes[i] = node.Hash()
			} else {
				hashes[i] = node.Raw()
			}
		}
	}

	hashes[16] = b.Value
	return hashes
}

func (b BranchNode) Hash() []byte {
	return crypto.Keccak256(b.Serialize())
}

func (b BranchNode) Serialize() []byte {
	return Serialize(b)
}

func (b BranchNode) HasValue() bool {
	return b.Value != nil
}

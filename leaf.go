package mptrie

import (
	"github.com/ethereum/go-ethereum/crypto"
)

type LeafNode struct {
	Path  []Nibble
	Value []byte
}

func (l LeafNode) Hash() []byte {
	return crypto.Keccak256(l.Serialize())
}

func (l LeafNode) Serialize() []byte {
	return Serialize(l)
}

func (l LeafNode) Raw() []interface{} {
	path := ToBytes(ToPrefixed(l.Path, true))
	raw := []interface{}{path, l.Value}
	return raw
}

func NewLeafNodeFromNibbles(nibbles []Nibble, value []byte) *LeafNode {
	return &LeafNode{
		Path:  nibbles,
		Value: value,
	}
}

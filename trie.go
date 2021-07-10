package main

import (
	"errors"
)

type Trie struct {
	root Node
}

func NewTrie() *Trie {
	return new(Trie)
}

func (t *Trie) Hash() []byte {
	if t.root == nil {
		return EmptyNodeHash
	}

	return t.root.Hash()
}

func (t *Trie) Put(key, value []byte) error {
	node := &t.root
	nibbles := FromBytes(key)

	if len(nibbles) <= 0 {
		return errors.New("cannot insert empty keys")
	}

	for {
		if *node == nil {
			leaf := NewLeafNodeFromNibbles(nibbles, value)
			*node = leaf
			return nil
		}

		if leaf, ok := (*node).(*LeafNode); ok {
			matched := PrefixMatchedLen(leaf.Path, nibbles)

			// all the leaf.Path matches with nibbles then update the value
			if matched == len(leaf.Path) && matched == len(nibbles) {
				newleaf := NewLeafNodeFromNibbles(leaf.Path, value)
				*node = newleaf
				return nil
			}

			branch := NewBranchNode()

			if matched == len(leaf.Path) {
				branch.SetValue(leaf.Value)
			}

			if matched == len(nibbles) {
				branch.SetValue(value)
			}

			// if there is matched nibbles, an extension node will be created
			if matched > 0 {
				ext := NewExtensionNode(leaf.Path[:matched], branch)
				*node = ext
			} else {
				*node = branch
			}

			if matched < len(leaf.Path) {
				branchNibble, leafNibbles := leaf.Path[matched], leaf.Path[matched+1:]
				newLeaf := NewLeafNodeFromNibbles(leafNibbles, leaf.Value)

				branch.SetBranch(branchNibble, newLeaf)
			}

			if matched < len(nibbles) {
				branchNibble, leafNode := nibbles[matched], nibbles[matched+1:]
				newLeaf := NewLeafNodeFromNibbles(leafNode, value)

				branch.SetBranch(branchNibble, newLeaf)
			}

			return nil
		}

		if branch, ok := (*node).(*BranchNode); ok {
			branchNibble, remaining := nibbles[0], nibbles[1:]
			nibbles = remaining
			node = &branch.Branches[int(branchNibble)]

			continue
		}

		if ext, ok := (*node).(*ExtensionNode); ok {
			matched := PrefixMatchedLen(ext.Path, nibbles)

			if matched < len(ext.Path) {
				extNibbles, branchNibble, extRemaining := ext.Path[:matched], ext.Path[matched], ext.Path[matched+1:]
				newBranchNibble, newLeafNibbles := nibbles[matched], nibbles[matched+1:]

				branch := NewBranchNode()
				if len(extRemaining) == 0 {
					branch.SetBranch(branchNibble, ext.Next)
				} else {
					newExt := NewExtensionNode(extRemaining, ext.Next)
					branch.SetBranch(branchNibble, newExt)
				}

				newleaf := NewLeafNodeFromNibbles(newLeafNibbles, value)
				branch.SetBranch(newBranchNibble, newleaf)

				if len(extNibbles) == 0 {
					*node = branch
				} else {
					*node = NewExtensionNode(extNibbles, branch)
				}

				return nil
			}

			nibbles = nibbles[matched:]
			node = &ext.Next
			continue
		}
	}
}

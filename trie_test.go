package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPutLeafDataIntoTrie(t *testing.T) {
	key := []byte("account.address")
	value := []byte("XYZABCDEF")

	trie := NewTrie()
	err := trie.Put(key, value)
	require.NoError(t, err)

	require.IsType(t, &LeafNode{}, trie.root)

	leaf := trie.root.(*LeafNode)
	expectedNibbles := FromBytes(key)
	require.Equal(t, value, leaf.Value)
	require.Equal(t, expectedNibbles, leaf.Path)
}

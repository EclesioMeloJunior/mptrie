package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPut_ShouldReturnLeafWhenTrieIsEmpty(t *testing.T) {
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

func TestPut_ShouldReturnExtensionNode(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("accounts.address"), []byte("some_fake_addresss"))
	require.NoError(t, err)
	require.IsType(t, &LeafNode{}, trie.root)

	err = trie.Put([]byte("accounts.value"), []byte("9000"))
	require.NoError(t, err)
	require.IsType(t, &ExtensionNode{}, trie.root)
}

func TestPut_ShouldReturnLeafWhenUpdateValue(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("accounts.address"), []byte("some_fake_addresss"))
	require.NoError(t, err)
	require.IsType(t, &LeafNode{}, trie.root)
	require.Equal(t, trie.root.(*LeafNode).Value, []byte("some_fake_addresss"))

	err = trie.Put([]byte("accounts.address"), []byte("another_address"))
	require.NoError(t, err)
	require.IsType(t, &LeafNode{}, trie.root)
	require.Equal(t, trie.root.(*LeafNode).Value, []byte("another_address"))
}

func TestPut_ShouldReturnBranchNodeWhenThereIsNoMatch(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("accounts.balance"), []byte("10000"))
	require.NoError(t, err)

	err = trie.Put([]byte("system.version"), []byte("1.0.0.0"))
	require.NoError(t, err)

	require.IsType(t, &BranchNode{}, trie.root)

	firstNibbles, secondNibbles := FromBytes([]byte("accounts.balance")), FromBytes([]byte("system.version"))
	branch := trie.root.(*BranchNode)

	require.NotNil(t, branch.Branches[int(firstNibbles[0])])
	require.NotNil(t, branch.Branches[int(secondNibbles[0])])

	firstLeaf := branch.Branches[int(firstNibbles[0])].(*LeafNode)
	secondLeaf := branch.Branches[int(secondNibbles[0])].(*LeafNode)

	require.Equal(t, firstNibbles[1:], firstLeaf.Path)
	require.Equal(t, secondNibbles[1:], secondLeaf.Path)

	require.Equal(t, firstLeaf.Value, []byte("10000"))
	require.Equal(t, secondLeaf.Value, []byte("1.0.0.0"))
}

func TestPut_ShouldStoreLeafValueAtBranchNode(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("transfer.input"), []byte("my-address"))
	require.NoError(t, err)

	err = trie.Put([]byte("transfer.input.value"), []byte("50"))
	require.NoError(t, err)

	require.IsType(t, &ExtensionNode{}, trie.root)

	expectedExtNibbles := FromBytes([]byte("transfer.input"))
	extnode := trie.root.(*ExtensionNode)

	require.Equal(t, expectedExtNibbles, extnode.Path)
	require.NotNil(t, extnode.Next)

	require.IsType(t, &BranchNode{}, extnode.Next)
	branch := extnode.Next.(*BranchNode)

	require.True(t, branch.HasValue())
	require.Equal(t, []byte("my-address"), branch.Value)
}

func TestPut_ShouldStoreNewValueAtBranchNode(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("transfer.input.value"), []byte("50"))
	require.NoError(t, err)

	err = trie.Put([]byte("transfer.input"), []byte("my-address"))
	require.NoError(t, err)

	require.IsType(t, &ExtensionNode{}, trie.root)

	expectedExtNibbles := FromBytes([]byte("transfer.input"))
	extnode := trie.root.(*ExtensionNode)

	require.Equal(t, expectedExtNibbles, extnode.Path)
	require.NotNil(t, extnode.Next)

	require.IsType(t, &BranchNode{}, extnode.Next)
	branch := extnode.Next.(*BranchNode)

	require.True(t, branch.HasValue())
	require.Equal(t, []byte("my-address"), branch.Value)
}

func TestPut_ShouldReturnErrWhenKeyEmpty(t *testing.T) {
	trie := NewTrie()
	err := trie.Put([]byte(""), []byte(""))

	require.Error(t, err)
}

func TestPut_WhenRootIsBranchNodeWithEmptySlot_ShouldAddLeafNode(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("zirst_info"), []byte("address"))
	require.NoError(t, err)

	err = trie.Put([]byte("other_info"), []byte("8000"))
	require.NoError(t, err)

	require.IsType(t, &BranchNode{}, trie.root)

	err = trie.Put([]byte("10hird_info"), []byte("some-hash"))
	require.NoError(t, err)

	require.IsType(t, &BranchNode{}, trie.root)
	branch := trie.root.(*BranchNode)

	firstNibbles := FromBytes([]byte("prevote"))
	secondNibbles := FromBytes([]byte("transfer.input"))
	thirdNibbles := FromBytes([]byte("block.header"))

	require.NotNil(t, branch.Branches[int(firstNibbles[0])])
	require.NotNil(t, branch.Branches[int(secondNibbles[0])])
	require.NotNil(t, branch.Branches[int(thirdNibbles[0])])

	require.IsType(t, &LeafNode{}, branch.Branches[int(firstNibbles[0])])
	require.IsType(t, &LeafNode{}, branch.Branches[int(secondNibbles[0])])
	require.IsType(t, &LeafNode{}, branch.Branches[int(thirdNibbles[0])])
}

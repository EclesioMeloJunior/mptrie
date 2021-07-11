package mptrie

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

func TestPut_WhenRootIsExtensionNodeShouldAddNewLeafNode(t *testing.T) {
	trie := NewTrie()

	err := trie.Put([]byte("block.header"), []byte("some_hash"))
	require.NoError(t, err)

	err = trie.Put([]byte("block.number"), []byte("1"))
	require.NoError(t, err)

	require.IsType(t, &ExtensionNode{}, trie.root)

	newNibbles := FromBytes([]byte("block.number"))
	leafNibbles := FromBytes([]byte("block.header"))

	matched := PrefixMatchedLen(newNibbles, leafNibbles)

	extnode := trie.root.(*ExtensionNode)
	require.Equal(t, leafNibbles[:matched], extnode.Path)
	require.IsType(t, &BranchNode{}, extnode.Next)

	err = trie.Put([]byte("transfer.input"), []byte("1000"))
	require.NoError(t, err)

	require.IsType(t, &BranchNode{}, trie.root)
	txNibbles := FromBytes([]byte("transfer.input"))

	branch := trie.root.(*BranchNode)
	txLeafNode := branch.Branches[int(txNibbles[0])]

	require.IsType(t, &LeafNode{}, txLeafNode)

	extNode := branch.Branches[int(extnode.Path[0])]
	require.IsType(t, &ExtensionNode{}, extNode)
}

func TestPut_WhenExtensionDoesntHaveRemaining(t *testing.T) {
	firstKey, firstValue := []byte("transfer.to"), []byte("some-addr")
	secondKey, secondValue := []byte("transfer.input"), []byte("some-value")
	thirdKey, thirdValue := []byte("transfer.gas"), []byte("some-fee")

	trie := NewTrie()

	err := trie.Put(firstKey, firstValue)

	require.NoError(t, err)
	require.IsType(t, &LeafNode{}, trie.root)

	err = trie.Put(secondKey, secondValue)

	require.NoError(t, err)
	require.IsType(t, &ExtensionNode{}, trie.root)

	err = trie.Put(thirdKey, thirdValue)

	require.NoError(t, err)
	require.IsType(t, &ExtensionNode{}, trie.root)
}

package mptrie

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"

	ethtrie "github.com/ethereum/go-ethereum/trie"
)

type proposal struct {
	ProposersSig [][]byte
	Document     string
}

func createFakeProposals() []proposal {
	ps := make([]proposal, 3)
	ps[0] = proposal{
		ProposersSig: [][]byte{
			[]byte("author 01"),
			[]byte("author 02"),
		},
		Document: "example 01",
	}

	ps[1] = proposal{
		ProposersSig: [][]byte{
			[]byte("author 03"),
		},
		Document: "example 02",
	}

	ps[2] = proposal{
		ProposersSig: [][]byte{
			[]byte("author 04"),
			[]byte("author 05"),
		},
		Document: "example 03 - with updates",
	}

	return ps
}

func TestCreateProof(t *testing.T) {
	trie := NewTrie()
	proposals := createFakeProposals()

	for i, p := range proposals {
		encoded, err := rlp.EncodeToBytes(p)
		require.NoError(t, err)

		key := make([]byte, 4)
		binary.LittleEndian.PutUint32(key, uint32(i))

		err = trie.Put(key, encoded)
		require.NoError(t, err)
	}

	key := make([]byte, 4)
	binary.BigEndian.PutUint32(key, uint32(0))

	m := NewInMemoryStorage()
	err := CreateProof(key, trie, m)
	require.NoError(t, err)

	v, err := ethtrie.VerifyProof(common.BytesToHash(trie.Hash()), key, m)
	require.NoError(t, err)

	var p proposal
	err = rlp.DecodeBytes(v, &p)
	require.NoError(t, err)

	fmt.Println(v, string(p.ProposersSig[0]))
}

package mptrie

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNibble_ShouldReturnTrue(t *testing.T) {
	bs := []byte{'a', 'b', 0xff, 3, 44}
	ns := FromBytes(bs)

	for _, n := range ns {
		require.True(t, IsNibble(byte(n)))
	}

}

func TestFromByte_ShouldReturnNibbles(t *testing.T) {
	b := byte(0xff) // represents 11111111 (255) when do b >> 4 should return 00001111 1 + 2 + 4 + 8 (15)
	expected := []Nibble{15, 15}

	ns := FromByte(b)

	require.Equal(t, expected, ns)
	require.Len(t, ns, 2)
}

func TestFromBytes_ShouldReturnNibbles(t *testing.T) {
	bs := []byte{0xff, 0xff, 0xff, 0xff}
	ns := FromBytes(bs)
	expected := []Nibble{15, 15, 15, 15, 15, 15, 15, 15}

	require.Equal(t, expected, ns)
	require.Len(t, ns, len(bs)*2)
}

func TestToBytes_ShouldReturnTheInitialBytes(t *testing.T) {
	b := byte(0xff)   // 11111111
	ns := FromByte(b) // []Nibble{15, 15}

	bs := ToBytes(ns) // 15 (00001111) << 4 = 11110000 (240) + 15 = 255 (0xff)

	require.Len(t, bs, 1)
	require.Equal(t, b, bs[0])
}

func TestToPrefixed_WhenIsLeafShouldReturnRightPrefix(t *testing.T) {
	b := byte(0xff)
	ns := FromByte(b)

	prefixedNibbles := ToPrefixed(ns, true)
	expected := []Nibble{2, 0, 15, 15}

	require.Len(t, prefixedNibbles, 4)
	require.Equal(t, expected, prefixedNibbles)
}

func TestPrefixMatchedLen(t *testing.T) {
	key1 := []byte{'a', 'a', 'c'}
	key2 := []byte{'a', 'a'}

	expectMatchedLen := 4 // [a, b, c] [a, b] == 2 * 2 (each byte = 2 nibbles)
	value := PrefixMatchedLen(FromBytes(key1), FromBytes(key2))

	require.Equal(t, expectMatchedLen, value)

	putKey := FromBytes([]byte("account.Balance"))
	putValue := make([]byte, 4)
	binary.LittleEndian.PutUint32(putValue, 100)

	leafn := NewLeafNodeFromNibbles(putKey, putValue)

	equalPrefixed := PrefixMatchedLen(leafn.Path, putKey)
	require.Equal(t, len(putKey), equalPrefixed)
	require.Equal(t, len(leafn.Path), equalPrefixed)
}

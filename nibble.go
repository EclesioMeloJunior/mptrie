package main

import "fmt"

type Nibble byte

func IsNibble(nibble byte) bool {
	n := int(nibble)
	return n >= 0 && n <= 16
}

func FromByte(b byte) []Nibble {
	return []Nibble{
		Nibble(byte(b >> 4)),
		Nibble(byte(b % 16)),
	}
}

func FromBytes(bs []byte) []Nibble {
	ns := make([]Nibble, 0, len(bs)*2)

	for _, b := range bs {
		ns = append(ns, FromByte(b)...)
	}

	return ns
}

// ToPrefixed add nibble prefix to a slice of nibbles to make its length even
// the prefix indicts whether a node is a leaf node.
func ToPrefixed(ns []Nibble, isLeafNode bool) []Nibble {
	var pf []Nibble

	if len(ns)%2 > 0 {
		pf = []Nibble{1}
	} else {
		pf = []Nibble{0, 0}
	}

	prefixed := make([]Nibble, 0, len(pf)+len(ns))
	prefixed = append(prefixed, pf...)

	for _, n := range ns {
		prefixed = append(prefixed, Nibble(n))
	}

	if isLeafNode {
		prefixed[0] += 2
	}

	return prefixed
}

func ToBytes(ns []Nibble) []byte {
	buf := make([]byte, 0, len(ns)/2)

	for i := 0; i < len(ns); i += 2 {
		b := byte(ns[i]<<4) + byte(ns[i+1])
		buf = append(buf, b)
	}

	fmt.Println(buf)

	return buf
}

func PrefixMatchedLen(n []Nibble, b []Nibble) (matched int) {
	for i := 0; i < len(n) && i < len(b); i++ {
		if n[i] == b[i] {
			matched++
		} else {
			break
		}
	}

	return
}

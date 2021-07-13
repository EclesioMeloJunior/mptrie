package mptrie

import "errors"

var (
	ErrWhileProof = errors.New("problem while verify proof")
)

func CreateProof(key []byte, t *Trie, r KVWriter) error {
	node := t.root
	nibbles := FromBytes(key)
	var nodes []Node

	for {
		if node == nil {
			return errors.New("node is empty")
		}

		if leaf, ok := node.(*LeafNode); ok {
			matched := PrefixMatchedLen(nibbles, leaf.Path)
			if len(leaf.Path) == matched && len(nibbles) == matched {
				nodes = append(nodes, leaf)
				break
			}

			return errors.New("key not found, cannot generate proof")
		}

		if branch, ok := node.(*BranchNode); ok {
			nodes = append(nodes, branch)

			if len(nibbles) == 0 {
				break
			}

			b, remaining := nibbles[0], nibbles[1:]
			nibbles = remaining
			node = branch.Branches[b]
			continue
		}

		if ext, ok := node.(*ExtensionNode); ok {
			matched := PrefixMatchedLen(nibbles, ext.Path)
			if matched < len(ext.Path) {
				return errors.New("extension path doest match, cannot generate proof")
			}

			nodes = append(nodes, ext)

			nibbles = nibbles[matched:]
			node = ext.Next
			continue
		}
	}

	for _, n := range nodes {
		if err := r.Put(Hash(n), Serialize(n)); err != nil {
			return err
		}
	}

	return nil
}

// func VerifyProof(root, key []byte, w KVReader) ([]byte, error) {
// 	nibbles := FromBytes(key)
// 	want := root

// 	for {
// 		b, err := w.Get(want)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if b == nil {
// 			return nil, ErrWhileProof
// 		}

// 	}
// }

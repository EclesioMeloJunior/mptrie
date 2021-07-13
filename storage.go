package mptrie

type KVWriter interface {
	Put([]byte, []byte) error
	Delete([]byte) error
}

type KVReader interface {
	Has([]byte) error
	Get([]byte) ([]byte, error)
}

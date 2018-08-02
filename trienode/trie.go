package trie

type trieFunctions interface {
	add() (bool, error)
	remove() (bool, error)
	find() bool
	New() *Trie
}

// Trie is the radix tree being built
type Trie struct {
	root  *Node
	words uint64
}

package trienode

import (
	"sync"
)

// Node is how a trie has it's internal node representation
type Node struct {
	letter   rune
	level    int
	metadata interface{}
	leaf     bool
	children map[rune]*Node
	locks    sync.RWMutex
}

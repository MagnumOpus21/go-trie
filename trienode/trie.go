package trienode

import (
	"fmt"
)

type trieFunctions interface {
	Add(element string)
	Remove(element string) (bool, error)
	Find(element string) bool
}

// Trie is the radix tree being built
type Trie struct {
	Root  *Node
	Words uint64
}

// New returns a new Trie
func New() *Trie {
	rootInit := new(Node)
	rootInit.children = make(map[rune]*Node)
	return &Trie{Words: 0, Root: rootInit}
}

// addChildNode returns a new child node if the current node is nil
func (n *Node) addChildNode(letter rune, metaData interface{}, leafCheck bool) *Node {
	return &Node{letter: letter, level: n.level + 1, metadata: metaData,
		leaf: !leafCheck, children: make(map[rune]*Node)}
}

// remove tries to delete a word from the Trie
func (n *Node) remove(word []rune) {
	if len(word) == 1 {
		fmt.Println("Hurr ", word, n)
		if len(n.children[word[0]].children) == 0 {
			delete(n.children, word[0])
		} else {
			n.children[word[0]].leaf = false
		}
		fmt.Println("Remove: ", n.children)
		return
	}
	n.children[word[0]].remove(word[1:])
}

// Find returns a bool value regarding the status of the current
// word being queried for
func (t *Trie) Find(word string) (status bool) {
	wordRuned := []rune(word)
	node := t.Root
	for _, letter := range wordRuned {
		fmt.Println(letter, node)
		if nodeChild, ok := node.children[letter]; ok {
			node = nodeChild
			status = node.leaf
		} else {
			return false
		}
	}
	return status
}

// Add is used to insert an element into the trie
// Add is always assumed to never fail
func (t *Trie) Add(word string, metaData interface{}) (found bool, err error) {
	// Use Find to make sure element is never present inside the Trie
	found = t.Find(word)
	if found {
		return found, fmt.Errorf("Word already exists in the Trie")
	}
	wordRuned := []rune(word)
	node := t.Root
	for index, letter := range wordRuned {
		if nodeChild, ok := node.children[letter]; ok {
			node = nodeChild
			if index == len(wordRuned)-1 {
				node.leaf = true
				node.metadata = metaData
			}
		} else {
			//Obtain a lock on this node to prevent other go-routines
			// from writing to this node
			node.locks.Lock()
			node.children[letter] = node.addChildNode(letter, metaData, (index < len(wordRuned)-1))
			node.locks.Unlock()
			node = node.children[letter]
		}
	}
	t.Words++
	return !found, nil
}

// Remove removes an element from the Trie
func (t *Trie) Remove(word string) (status bool, err error) {
	if exists := t.Find(word); !exists {
		return status, fmt.Errorf("No such word in the trie")
	}
	wordRune := []rune(word)
	node := t.Root
	node.remove(wordRune)
	t.Words--
	return true, nil
}

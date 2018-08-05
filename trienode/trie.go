package trienode

import (
	"fmt"
	"sort"
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

//Words is a string array
type Words []string

// Sort interface
func (w Words) Len() int {
	return len(w)
}

func (w Words) Swap(a, b int) {
	w[a], w[b] = w[b], w[a]
}

func (w Words) Less(a, b int) bool {
	if len(w[a]) >= len(w[b]) {
		return w[a] <= w[b]
	}
	return w[a] > w[b]
}

// New returns a new Trie
func New() *Trie {
	rootInit := new(Node)
	rootInit.children = make(map[rune]*Node)
	return &Trie{Words: 0, Root: rootInit}
}

// addChildNode returns a new child node if the current node is nil
func (n *Node) addChildNode(letter rune, metaData interface{}, leafCheck bool) *Node {
	if leafCheck == true {
		metaData = nil
	}
	return &Node{letter: letter, level: n.level + 1, metadata: metaData,
		leaf: !leafCheck, children: make(map[rune]*Node)}
}

func (n *Node) prefixSearch(prefix []rune) (matches Words) {
	node := n
	prefixFound := make([]rune, 0)
	for _, letter := range prefix {
		if nodeCheck, ok := node.children[letter]; ok {
			node = nodeCheck
			prefixFound = append(prefixFound, letter)
		} else {
			break
		}
	}
	matches = node.prefixSearchHelper(prefix, matches)
	sort.Sort(matches)
	return matches
}

func (n *Node) prefixSearchHelper(prefix []rune, match []string) []string {
	if n.leaf {
		match = append(match, string(prefix))
	}
	for _, candidate := range n.children {
		newPrefix := append(prefix, candidate.letter)
		tempMatch := candidate.prefixSearchHelper(newPrefix, make([]string, 0))
		for _, word := range tempMatch {
			match = append(match, word)
		}
	}
	return match
}

// remove tries to delete a word from the Trie
func (n *Node) remove(word []rune) {
	if len(word) == 1 {
		if len(n.children[word[0]].children) == 0 {
			delete(n.children, word[0])
		} else {
			n.children[word[0]].leaf = false
		}
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
				node.locks.Lock()
				node.leaf = true
				node.metadata = metaData
				node.locks.Unlock()
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

// PrefixSearch returns all the words that have the given prefix
func (t *Trie) PrefixSearch(prefix string) (matches []string) {
	matches = t.Root.prefixSearch([]rune(prefix))
	return matches
}

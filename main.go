package main

import (
	"fmt"

	trie "github.com/go-trie/trienode"
)

func main() {
	var trieData *trie.Trie
	trieData = trie.New()
	if _, err := trieData.Add("Hello", "Coming through"); err != nil {
		fmt.Println(err)
	}
	if _, err := trieData.Add("Hell", "???"); err != nil {
		fmt.Println(err)
	}
	if _, err := trieData.Add("Jell", "Testing"); err != nil {
		fmt.Println(err)
	}
	if _, err := trieData.Add("Jello", "Testing add"); err != nil {
		fmt.Println(err)
	}
	if _, err := trieData.Add("Help", "Testing add"); err != nil {
		fmt.Println(err)
	}
	if _, err := trieData.Add("Hell", "Testing add"); err != nil {
		fmt.Println(err, "Hell")
	}
	fmt.Println("Words in the Trie: ", trieData.Words)
	if _, err := trieData.Remove("Hell"); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Remove check")
	if ok := trieData.Find(("Hell")); ok != false {
		fmt.Println("Remove doesn't work")
	}
	trieData.Remove("Help")
	fmt.Println("Words in the Trie: ", trieData.Words)
}

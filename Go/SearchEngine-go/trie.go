package main

type TrieNode struct {
	Children map[rune]*TrieNode
	IsEnd    bool
	Words    map[string]bool
}

type Trie struct {
	Root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{Root: &TrieNode{Children: make(map[rune]*TrieNode), Words: make(map[string]bool)}}
}

func (t *Trie) Insert(word, docID string) {
	node := t.Root
	for _, char := range word {
		if node.Children[char] == nil {
			node.Children[char] = &TrieNode{Children: make(map[rune]*TrieNode), Words: make(map[string]bool)}
		}
		node = node.Children[char]
		node.Words[docID] = true
	}
	node.IsEnd = true
}

func (t *Trie) Suggest(prefix string, limit int) []string {
	node := t.Root
	for _, char := range prefix {
		if node.Children[char] == nil {
			return nil
		}
		node = node.Children[char]
	}
	var docs []string
	for docID := range node.Words {
		docs = append(docs, docID)
		if len(docs) >= limit {
			break
		}
	}
	return docs
}

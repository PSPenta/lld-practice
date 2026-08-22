package main

type Posting struct {
	DF   int
	Docs map[string]int
}

type InvertedIndex struct {
	Index    map[string]*Posting
	DocCount int
}

func NewInvertedIndex() *InvertedIndex {
	return &InvertedIndex{Index: make(map[string]*Posting)}
}

func (i *InvertedIndex) AddDocument(docID string, tokens []string) {
	i.DocCount++
	freq := make(map[string]int)
	for _, token := range tokens {
		freq[token]++
	}
	for token, count := range freq {
		if i.Index[token] == nil {
			i.Index[token] = &Posting{Docs: make(map[string]int)}
		}
		i.Index[token].DF++
		i.Index[token].Docs[docID] = count
	}
}

func (i *InvertedIndex) GetPosting(token string) *Posting {
	return i.Index[token]
}

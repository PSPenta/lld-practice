package main

type SearchResult struct {
	DocID string
	Text  string
}

type SearchEngine struct {
	tokenizer *Tokenizer
	trie      *Trie
	index     *InvertedIndex
	ranker    *Ranker
	documents map[string]string
}

func NewSearchEngine() *SearchEngine {
	index := NewInvertedIndex()
	return &SearchEngine{
		tokenizer: NewTokenizer(nil, 2),
		trie:      NewTrie(),
		index:     index,
		ranker:    NewRanker(index),
		documents: make(map[string]string),
	}
}

func (s *SearchEngine) AddDocument(docID, text string) {
	if docID == "" || text == "" {
		return
	}
	s.documents[docID] = text
	tokens := s.tokenizer.Tokenize(text)
	s.index.AddDocument(docID, tokens)
	for _, token := range tokens {
		s.trie.Insert(token, docID)
	}
}

func (s *SearchEngine) Search(query string, limit int) []SearchResult {
	if limit == 0 {
		limit = 10
	}
	tokens := s.tokenizer.Tokenize(query)
	rankedDocs := s.ranker.Score(tokens)
	if len(rankedDocs) > limit {
		rankedDocs = rankedDocs[:limit]
	}
	var results []SearchResult
	for _, docID := range rankedDocs {
		results = append(results, SearchResult{DocID: docID, Text: s.documents[docID]})
	}
	return results
}

func (s *SearchEngine) Autocomplete(prefix string, limit int) []string {
	if limit == 0 {
		limit = 5
	}
	docs := s.trie.Suggest(prefix, limit)
	var texts []string
	for _, docID := range docs {
		texts = append(texts, s.documents[docID])
	}
	return texts
}

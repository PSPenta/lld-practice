package main

import (
	"math"
	"sort"
)

type Ranker struct {
	Index *InvertedIndex
}

func NewRanker(index *InvertedIndex) *Ranker {
	return &Ranker{Index: index}
}

func (r *Ranker) Score(queryTokens []string) []string {
	scores := make(map[string]float64)
	for _, token := range queryTokens {
		posting := r.Index.GetPosting(token)
		if posting == nil {
			continue
		}
		idf := math.Log(float64(1+r.Index.DocCount)/float64(1+posting.DF)) + 1
		for docID, tf := range posting.Docs {
			scores[docID] += float64(tf) * idf
		}
	}
	type scoredDoc struct {
		docID string
		score float64
	}
	var ranked []scoredDoc
	for docID, score := range scores {
		ranked = append(ranked, scoredDoc{docID, score})
	}
	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].score > ranked[j].score
	})
	var docIDs []string
	for _, item := range ranked {
		docIDs = append(docIDs, item.docID)
	}
	return docIDs
}

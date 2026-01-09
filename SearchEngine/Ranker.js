class Ranker {
  constructor(index) {
    this.index = index;
  }

  score(queryTokens) {
    const scores = {};

    for (const token of queryTokens) {
      const posting = this.index.getPosting(token);
      if (!posting) continue;

      const { df, docs } = posting;
      const idf = Math.log((1 + this.index.docCount) / (1 + df)) + 1;

      for (const [docId, tf] of Object.entries(docs)) {
        const tfidf = tf * idf;
        scores[docId] = (scores[docId] || 0) + tfidf;
      }
    }

    // return sorted doc list
    return Object.entries(scores)
      .sort((a, b) => b[1] - a[1])
      .map(([docId]) => docId);
  }
}

module.exports = { Ranker };

class InvertedIndex {
  constructor() {
    this.index = Object.create(null);
    this.docCount = 0;
  }

  addDocument(docId, tokens) {
    this.docCount++;

    const freqMap = new Map();
    for (const token of tokens) {
      freqMap.set(token, (freqMap.get(token) || 0) + 1);
    }

    for (const [token, count] of freqMap.entries()) {
      if (!this.index[token]) {
        this.index[token] = { df: 0, docs: {} };
      }
      this.index[token].df++;
      this.index[token].docs[docId] = count;
    }
  }

  getPosting(token) {
    return this.index[token] || null;
  }
}

module.exports = { InvertedIndex };

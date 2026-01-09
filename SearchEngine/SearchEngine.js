const { Tokenizer } = require("./Tokenizer");
const { Trie } = require("./Trie");
const { InvertedIndex } = require("./InvertedIndex");
const { Ranker } = require("./Ranker");

class SearchEngine {
  constructor() {
    this.tokenizer = new Tokenizer();
    this.trie = new Trie();
    this.index = new InvertedIndex();
    this.ranker = new Ranker(this.index);

    this.documents = {};  // docId → raw text
  }

  addDocument(docId, text) {
    if (!docId || !text) return;

    this.documents[docId] = text;

    const tokens = this.tokenizer.tokenize(text);

    // add to inverted index
    this.index.addDocument(docId, tokens);

    // add tokens to Trie for autocomplete
    for (const token of tokens) {
      this.trie.insert(token, docId);
    }
  }

  search(query, options = {}) {
    const tokens = this.tokenizer.tokenize(query);

    const rankedDocs = this.ranker.score(tokens).slice(0, options.limit || 10);

    return rankedDocs.map(docId => ({
      docId,
      text: this.documents[docId]
    }));
  }

  autocomplete(prefix, limit = 5) {
    const docs = this.trie.suggest(prefix.toLowerCase(), limit);
    return docs.map(docId => this.documents[docId]);
  }
}

module.exports = { SearchEngine };

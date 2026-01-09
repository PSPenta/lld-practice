class TrieNode {
  constructor() {
    this.children = Object.create(null);
    this.isEnd = false;
    this.words = new Set();  // limit suggestions by storing docs
  }
}

class Trie {
  constructor() {
    this.root = new TrieNode();
  }

  insert(word, docId) {
    let node = this.root;

    for (const char of word) {
      if (!node.children[char]) {
        node.children[char] = new TrieNode();
      }
      node = node.children[char];
      node.words.add(docId);
    }

    node.isEnd = true;
  }

  // returns top N ranked suggestions
  suggest(prefix, limit = 10) {
    let node = this.root;

    for (const char of prefix.toLowerCase()) {
      if (!node.children[char]) return [];
      node = node.children[char];
    }

    return Array.from(node.words).slice(0, limit);
  }
}

module.exports = { Trie };

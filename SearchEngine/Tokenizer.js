class Tokenizer {
  constructor(options = {}) {
    this.stopWords = new Set(options.stopWords || [
      "the", "is", "and", "or", "to", "in", "on", "for", "of", "a", "an"
    ]);

    this.minTokenLength = options.minTokenLength || 2;
  }

  tokenize(text = "") {
    if (!text || typeof text !== "string") return [];

    return text
      .toLowerCase()
      .replace(/[^\p{L}\p{N}\s]+/gu, "")   // remove punctuation safely
      .split(/\s+/)
      .filter(t => t.length >= this.minTokenLength)
      .filter(t => !this.stopWords.has(t));
  }
}

module.exports = { Tokenizer };

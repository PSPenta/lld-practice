const { SearchEngine } = require('./SearchEngine');

const engine = new SearchEngine();

engine.addDocument("1", "Apple launches new MacBook Pro with M4 chip");
engine.addDocument("2", "Samsung introduces AI-powered Galaxy devices");
engine.addDocument("3", "Apple releases major update for iOS 19");

console.log("🔍 Search: 'apple m4'");
console.log(engine.search("apple m4"));

console.log("🔍 Search: 'apple 19'");
console.log(engine.search("apple 19"));

console.log("💡 Autocomplete: 'mac'");
console.log(engine.autocomplete("mac"));

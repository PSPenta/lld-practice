class Node {
  constructor(key, value, ttl) {
    this.key = key;
    this.value = value;
    this.ttl = ttl;
    this.next = null;
    this.prev = null;
  }
}

class Redis {
  constructor(capacity) {
    this.allNodes = new Map();
    this.size = 0;
    this.capacity = capacity;
    this.head = null;
    this.tail = null;
    this.defaultTTL = 3600;      // 1 hour in seconds
  }

  set(key, value, expiry = this.defaultTTL) {
    let node = new Node(key, value, Date.now() + (expiry * 1000));
    this.#removeNode(this.allNodes.get(key));
    this.#addNode(node);

    return true;
  }

  get(key) {
    let node = this.allNodes.get(key);

    if (!node) return null;

    if (node.ttl > Date.now()) {
      this.#removeNode(node);
      this.#addNode(node);
    } else {
      this.#removeNode(node);
      return null;
    }

    return node.value;
  }

  #addNode(node) {
    if (!this.head) {
      this.head = node;
      this.tail = node;
      this.size++;
      this.allNodes.set(node.key, node);
      return true;
    }

    node.next = this.head;
    this.head.prev = node;
    this.head = node;
    this.size++;
    this.allNodes.set(node.key, node);

    if (this.size > this.capacity) {
      this.#evict();
    }
  }

  #removeNode(node) {
    if (!node) return;

    let prev = node.prev;
    let next = node.next;

    if (prev) {
      prev.next = next;
    } else {
      this.head = next;
    }

    if (next) {
      next.prev = prev;
    } else {
      this.tail = prev;
    }

    this.allNodes.delete(node.key);
    this.size--;
  }

  #evict() {
    let expiredNode = this.#getExpiredNode();
    if (expiredNode) {
      this.#removeNode(expiredNode);
      return true;
    }

    this.#removeNode(this.tail);
    return true;
  }

  #getExpiredNode() {
    let currentTime = Date.now();
    let node = this.tail;
    while (node && node.ttl > currentTime) {
      node = node.prev;
    }

    if (node && node.ttl <= currentTime) return node;

    return null;
  }
}

let redis = new Redis(10);
redis.set('test 1', 1);
console.log(redis.get('test 1'));
redis.set('test 2', 2);
console.log(redis.get('test 2'));
redis.set('test 3', 3);
console.log(redis.get('test 3'));
redis.set('test 4', 4);
console.log(redis.get('test 4'));
redis.set('test 5', 5);
console.log(redis.get('test 5'));
redis.set('test 6', 6);
console.log(redis.get('test 6'));
redis.set('test 7', 7);
console.log(redis.get('test 7'));
redis.set('test 8', 8);
console.log(redis.get('test 8'));
redis.set('test 9', 9);
console.log(redis.get('test 9'));
redis.set('test 10', 10);
console.log(redis.get('test 10'));
redis.set('test 11', 11);
console.log(redis.get('test 11'));
redis.set('test 12', 12);
console.log(redis.get('test 12'));
redis.set('test 13', 13);
console.log(redis.get('test 13'));
redis.set('test 14', 14);
console.log(redis.get('test 14'));
redis.set('test 15', 15);
console.log(redis.get('test 15'));
redis.set('test 16', 16);

console.log(redis.get('test 16'));
console.log(redis.get('test 1'));
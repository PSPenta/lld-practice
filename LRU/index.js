class Node {
  constructor(key, value) {
    this.key = key;
    this.value = value;
    this.next = null;
    this.prev = null;
  }
}

class LRUCache {
  constructor(capacity) {
    this.capacity = capacity;
    this.map = new Map();
    this.head = null;
    this.tail = null;
  }

  get(key) {
    const node = this.map.get(key);
    if(!node) return -1;

    this.#removeNode(node);
    this.#addNode(node);

    return node.value;
  }

  set(key, value) {
    const existingNode = this.map.get(key);

    if(existingNode) {
      existingNode.value = value;
      if (existingNode !== this.head) {
        this.#removeNode(existingNode); // Remove the node from its current position
        this.#addNode(existingNode); // Add the node to the front of the list
      }
    } else {
      const node = new Node(key, value);
      this.map.set(key, node);
      this.#addNode(node);
      if(this.map.size > this.capacity) {
        let tail = this.tail;
        this.map.delete(tail);
        this.#removeNode(tail);
      }
    }

    return 1;
  }


  // Helper method to remove a node from the doubly linked list
  #removeNode(node) {
    // If the node is the head, update the head pointer
    if (node === this.head) {
      this.head = node.next;
    }

    // If the node is the tail, update the tail pointer
    if (node === this.tail) {
      this.tail = node.prev;
    }

    // If the node has a previous node, update its next pointer
    if (node.prev) {
      node.prev.next = node.next;
    }

    // If the node has a next node, update its prev pointer
    if (node.next) {
      node.next.prev = node.prev;
    }
  }

  // Helper method to add a node to the front of the doubly linked list
  #addNode(node) {
    // If the list is empty, set the node as the head and the tail
    if (!this.head) {
      this.head = node;
      this.tail = node;
    } else {
      // If the list is not empty, insert the node before the head and update the head pointer
      node.next = this.head;
      this.head.prev = node;
      this.head = node;
    }
  }
}



const lru = new LRUCache(3);

lru.set(1, 'A');

lru.set(2, 'B');

lru.set(3, 'C');

lru.set(4, 'D');

console.log(lru.get(1));

console.log(lru.get(4));

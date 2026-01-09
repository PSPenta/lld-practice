class FIFOQueue {
  constructor(capacity) {
    this.capacity = capacity;
    this.front = -1;
    this.rear = -1;
    this.queue = new Array(capacity).fill().map(_ => null);
  }

  enqueue(value) {
    const rear = (this.rear + 1) % this.capacity;
    if(this.queue[rear] !== null) {
      console.log('queue is full');
      return -1;
    }
    this.queue[rear] = value;
    this.rear = rear;
  }

  dequeue() {
    const front = (this.front + 1) % this.capacity;
    if(this.queue[front] === null) {
      console.log('queue is empty');
      return -1;
    }
    const value = this.queue[front];
    this.queue[front] = null;
    this.front = front;
    return value;
  }
}

const fifo = new FIFOQueue(5);
fifo.enqueue(1);
fifo.enqueue(2);
fifo.enqueue(3);
fifo.enqueue(4);
fifo.enqueue(5);
// console.log(fifo.dequeue());
// console.log(fifo.dequeue());
// console.log(fifo.dequeue());
// console.log(fifo.dequeue());
fifo.enqueue(6);






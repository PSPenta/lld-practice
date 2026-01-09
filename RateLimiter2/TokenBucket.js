class TokenBucket {
  constructor(capacity, refillRate) {
    this.capacity = capacity;
    this.refillRate = refillRate;
    this.bucket = new Map();
  }

  refill(ip) {
    const currentTime = Date.now();
    const ipRequests = this.bucket.get(ip);

    const timeElapsed = (currentTime - ipRequests.lastRefilledAt) / 1000;
    const refills = Math.floor(timeElapsed * this.refillRate);

    if (refills > 0) {
      this.bucket.set(ip, {
        tokens: Math.min(this.capacity, ipRequests.tokens + refills),
        lastRefilledAt: currentTime
      });
    }
  }

  isAllowed(ip) {
    if (!this.capacity) return false;

    if (!this.bucket.has(ip)) {
      this.bucket.set(ip, {
        tokens: this.capacity - 1,
        lastRefilledAt: Date.now()
      });

      return true;
    }

    if (this.refillRate) {
      this.refill(ip);
    }

    const ipRequests = this.bucket.get(ip);
    if (ipRequests.tokens > 0) {
      this.bucket.set(ip, {
        tokens: ipRequests.tokens - 1,
        lastRefilledAt: ipRequests.lastRefilledAt
      });
      return true;
    }

    return false;
  }
}

const limiter = new TokenBucket(0, 0);
const ip = "1.2.3.4";

let i = 1;
const id = setInterval(() => {
  if (i > 20) clearInterval(id);
  console.log(`Request ${i}:`, limiter.isAllowed(ip));
  i++;
}, 500);

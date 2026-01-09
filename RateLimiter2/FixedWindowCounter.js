class FixedWindowCounter {
  limit = 0;
  windowMs = 0;
  requests = new Map();
  windowStart = Date.now();

  constructor(limit, windowMs) {
    this.limit = limit;
    this.windowMs = windowMs;
    this.windowStart = Date.now();
  }

  isAllowed(ip) {
    let currentTime = Date.now();
    if (currentTime - this.windowStart >= this.windowMs) {
      this.windowStart = currentTime;
      this.requests = new Map();
    }

    if (!this.requests.has(ip)) {
      this.requests.set(ip, 1);
      return true;
    }

    let currentCount = this.requests.get(ip);
    if (currentCount >= this.limit) {
      return false;
    }

    this.requests.set(ip, currentCount + 1);
    return true;
  }
};

const limiter = new FixedWindowCounter(3, 10000); // 3 req / 10s
const ip = "1.2.3.4";

for (let i = 1; i <= 5; i++) {
  console.log(`Request ${i}:`, limiter.isAllowed(ip));
}

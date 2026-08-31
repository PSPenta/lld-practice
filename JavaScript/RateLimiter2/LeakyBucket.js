const RateLimiterStrategy = require('./RateLimiterStrategy');

class LeakyBucket extends RateLimiterStrategy {
  constructor(capacity, leakRate) {
    super();
    this.capacity = capacity;
    this.leakRate = leakRate;
    this.bucket = new Map();
  }

  leak(ip) {
    const currentTime = Date.now();
    const ipRequests = this.bucket.get(ip);

    const timeElapsed = (currentTime - ipRequests.lastLeakedAt) / 1000;
    const leaks = Math.floor(this.leakRate * timeElapsed);

    if (leaks > 0) {
      this.bucket.set(ip, {
        totalRequests: Math.max(0, ipRequests.totalRequests - leaks),
        lastLeakedAt: currentTime,
      });
    }
  }

  isAllowed(ip) {
    if (!this.capacity) return false;

    if (!this.bucket.has(ip)) {
      this.bucket.set(ip, {
        totalRequests: 1,
        lastLeakedAt: Date.now(),
      });

      return true;
    }

    if (this.leakRate) {
      this.leak(ip);
    }

    const ipRequests = this.bucket.get(ip);
    if (ipRequests.totalRequests < this.capacity) {
      this.bucket.set(ip, {
        totalRequests: ipRequests.totalRequests + 1,
        lastLeakedAt: ipRequests.lastLeakedAt,
      });
      return true;
    }

    return false;
  }
}

module.exports = LeakyBucket;

const RateLimiterStrategy = require('./RateLimiterStrategy');

class FixedWindowCounter extends RateLimiterStrategy {
  limit = 0;
  windowMs = 0;
  requests = new Map();
  windowStart = Date.now();

  constructor(limit, windowMs) {
    super();
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
}

module.exports = FixedWindowCounter;

const RateLimiterStrategy = require('./RateLimiterStrategy');

class SlidingWindowLog extends RateLimiterStrategy {
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
    }

    if (!this.requests.has(ip)) {
      this.requests.set(ip, [{ time: currentTime }]);
      return true;
    }

    let oldRequests = this.requests.get(ip);
    oldRequests = oldRequests.filter(
      (req) => currentTime - req.time < this.windowMs,
    );

    if (oldRequests.length < this.limit) {
      this.requests.set(ip, [...oldRequests, { time: currentTime }]);
      return true;
    }

    return false;
  }
}

module.exports = SlidingWindowLog;

class RateLimiter {
  constructor(strategy) {
    this.strategy = strategy;
  }

  isAllowed(ip) {
    if (this.strategy) {
      return this.strategy.isAllowed(ip);
    }

    throw new Error("Please set a strategy!");
  }
}

module.exports = RateLimiter;

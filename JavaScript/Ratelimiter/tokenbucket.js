const RATE_LIMIT = 5;
const REFILL_RATE_PER_SECOND = 1; // Tokens added per second
const userBuckets = new Map();

class TokenBucket {
  constructor(limit, refillRate) {
    this.limit = limit;
    this.tokens = limit; // Start with the bucket full
    this.refillRate = refillRate;
    this.lastRefill = Date.now(); // Last time tokens were refilled
  }

  refill() {
    const now = Date.now();
    const timeElapsed = (now - this.lastRefill) / 1000;
    const tokensToAdd = Math.floor(timeElapsed * this.refillRate);

    this.tokens = Math.min(this.limit, this.tokens + tokensToAdd);
    this.lastRefill = now;
  }

  allowRequest() {
    this.refill();
    if (this.tokens > 0) {
      this.tokens--;
      return true;
    }
    return false;
  }
}
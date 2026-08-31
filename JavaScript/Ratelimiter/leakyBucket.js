const RATE_LIMIT = 5;
const LEAK_RATE_IN_MS = 1000; // Leak 1 request per second
const userBuckets = new Map();

class LeakyBucket {
  constructor(limit, leakRate) {
    this.limit = limit;
    this.leakRate = leakRate;
    this.requests = 0;
    this.lastLeak = Date.now();
  }

  leak() {
    const now = Date.now();
    const timeElapsed = now - this.lastLeak;

    const leakedRequests = Math.floor(timeElapsed / this.leakRate);

    this.requests = Math.max(0, this.requests - leakedRequests);
    this.lastLeak = now;
  }

  allowRequest() {
    this.leak();
    if (this.requests < this.limit) {
      this.requests++;
      return true;
    }
    return false;
  }
}

// Leaky Bucket Middleware
const leakyBucketRateLimiter = (req, res, next) => {
  const userIP = req.ip;

  if (!userBuckets.has(userIP)) {
    userBuckets.set(userIP, new LeakyBucket(RATE_LIMIT, LEAK_RATE_IN_MS));
  }

  const bucket = userBuckets.get(userIP);

  if (bucket.allowRequest()) {
    next();
  } else {
    res
      .status(429)
      .json({ message: 'Too many requests. Please try again later.' });
  }
};

module.exports = leakyBucketRateLimiter;

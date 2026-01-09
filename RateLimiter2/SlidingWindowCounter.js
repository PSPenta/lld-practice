class SlidingWindowCounter {
  constructor(limit, windowMs) {
    this.limit = limit;
    this.windowMs = windowMs;
    this.requests = new Map();
  }

  isAllowed(ip) {
    const currentTime = Date.now();

    if (!this.requests.has(ip)) {
      this.requests.set(ip, {
        previousStart: currentTime - this.windowMs,
        previousCount: 0,
        currentStart: currentTime,
        currentCount: 0
      });
    }

    const window = this.requests.get(ip);

    // Let's assume current window started 4 seconds ago, so timeLapsed = 4
    let timeLapsed = currentTime - window.currentStart;

    // If timeLapsed >= windowMs, time to reset the window
    if (timeLapsed >= this.windowMs) {
      window.previousStart = window.currentStart;
      window.previousCount = window.currentCount;
      window.currentStart = currentTime;
      window.currentCount = 0;

      if (timeLapsed >= (this.windowMs * 2)) {
        window.previousCount = 0;
      }
      this.requests.set(ip, window);

      // Fresh window started
      timeLapsed = 0;
    }

    // 3 / 8 = ~0.4 (~40% time lapsed) ---- Ratio: 0-0.9
    const ratio = timeLapsed / this.windowMs;

    // 1 - 0.4 = 0.6 (skip ~40% requests, take ~60% of previousCount to take 1 entire window i.e. ~60% of previous and ~40% of current)
    const effectivePreviousCount = window.previousCount * (1 - ratio);
    const effectiveTotal = window.currentCount + effectivePreviousCount;
    if (effectiveTotal < this.limit) {
      window.currentCount++;
      return true;
    }

    return false;
  }
}

const limiter = new SlidingWindowCounter(8, 10000);
const ip = "1.2.3.4";

let i = 1;
const id = setInterval(() => {
  if (i > 20) clearInterval(id);
  console.log(`Request ${i}:`, limiter.isAllowed(ip));
  i++;
}, 1000);

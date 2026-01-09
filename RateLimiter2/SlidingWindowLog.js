class SlidingWindowLog {
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
    }

    if (!this.requests.has(ip)) {
      this.requests.set(ip, [{ time: currentTime }]);
      return true;
    }

    let oldRequests = this.requests.get(ip);
    oldRequests = oldRequests.filter(req => currentTime - req.time < this.windowMs);

    if (oldRequests.length < this.limit) {
      this.requests.set(ip, [ ...oldRequests, { time: currentTime } ]);
      return true;
    }

    return false;
  }
}

const limiter = new SlidingWindowLog(8, 10000);
const ip = "1.2.3.4";

let i = 1;
const id = setInterval(() => {
  if (i > 20) clearInterval(id);
  console.log(`Request ${i}:`, limiter.isAllowed(ip));
  i++;
}, 1000);

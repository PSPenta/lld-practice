class FixedWindowRatelimiter {
  constructor(windowMs, rateLimit) {
    this.windowMs = windowMs;
    this.rateLimit = rateLimit;
    this.requestLog = new Map(); // { userId, <List>requests }
  }

  isRequestAllowed(userId, path) {
    const currentTime = Date.now()

    const existingRequest = this.requestLog.get(userId);

    if(!existingRequest) {
      this.requestLog.set(userId, [{path, timestamp: currentTime}])
      return true;
    }

    const validRequests = this.requestLog.get(userId).filter(({timestamp}) => currentTime - timestamp < this.windowMs)

    if(validRequests.length >= this.rateLimit) return false;

    this.requestLog.set(userId, [...validRequests, { path, timestamp: currentTime }]);

    return true;
  }
}
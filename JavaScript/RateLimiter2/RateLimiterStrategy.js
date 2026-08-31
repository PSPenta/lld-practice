class RateLimiterStrategy {
  isAllowed(ip) {
    throw new Error("This is an abstract class!");
  }
}

module.exports = RateLimiterStrategy;

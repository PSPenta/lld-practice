/**
 * @interface RateLimiterStrategy (abstract base class)
 *
 * JavaScript has no native `interface` keyword. This class is the contract that
 * every rate-limiting algorithm must satisfy (Strategy pattern).
 *
 * @method isAllowed(ip) → boolean
 *
 * @see TokenBucket, FixedWindowCounter, SlidingWindowLog, SlidingWindowCounter, LeakyBucket
 */
class RateLimiterStrategy {
  isAllowed(ip) {
    throw new Error("This is an abstract class!");
  }
}

module.exports = RateLimiterStrategy;

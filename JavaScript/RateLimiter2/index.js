const RateLimiter = require("./RateLimiter");
const TokenBucket = require("./TokenBucket");
const FixedWindowCounter = require("./FixedWindowCounter");
const SlidingWindowLog = require("./SlidingWindowLog");
const SlidingWindowCounter = require('./SlidingWindowCounter');
const LeakyBucket = require("./LeakyBucket");

const rateLimiter = new RateLimiter(new TokenBucket(5, 1));
rateLimiter.isAllowed("127.0.0.1");

const rateLimiter2 = new RateLimiter(new FixedWindowCounter(5, 1000));
rateLimiter2.isAllowed("127.0.0.1");

const rateLimiter3 = new RateLimiter(new SlidingWindowLog(5, 1000));
rateLimiter3.isAllowed("127.0.0.1");

const rateLimiter4 = new RateLimiter(new LeakyBucket(5, 1));
rateLimiter4.isAllowed("127.0.0.1");

const rateLimiter5 = new RateLimiter(new SlidingWindowCounter(5, 1000));
rateLimiter5.isAllowed('127.0.0.1');

console.log(rateLimiter.isAllowed("127.0.0.1"));
console.log(rateLimiter2.isAllowed("127.0.0.1"));
console.log(rateLimiter3.isAllowed("127.0.0.1"));
console.log(rateLimiter4.isAllowed("127.0.0.1"));
console.log(rateLimiter5.isAllowed('127.0.0.1'));

console.log(rateLimiter.isAllowed("127.0.0.1"));
console.log(rateLimiter2.isAllowed("127.0.0.1"));
console.log(rateLimiter3.isAllowed("127.0.0.1"));
console.log(rateLimiter4.isAllowed("127.0.0.1"));
console.log(rateLimiter5.isAllowed('127.0.0.1'));

console.log(rateLimiter.isAllowed("127.0.0.1"));
console.log(rateLimiter2.isAllowed("127.0.0.1"));
console.log(rateLimiter3.isAllowed("127.0.0.1"));
console.log(rateLimiter4.isAllowed("127.0.0.1"));
console.log(rateLimiter5.isAllowed('127.0.0.1'));

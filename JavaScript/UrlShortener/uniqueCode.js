const BASE_62 = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
const DEFAULT_LENGTH = 6;
const SECRET = 9765439;               // A secret salt number to break the unpredictability

exports.generateShortCode = (db) => {
  let shortCode = '';

  // Can replace this with redis.incr('UNIQUE_COUNTER') to make atomic increment, so counter can remain unique across all servers.
  const counter = (db.get('UNIQUE_COUNTER') || 0) + 1;
  db.set('UNIQUE_COUNTER', counter);

  let temp = counter ^ SECRET;          // We can use Knuth multiplicative hash instead of this XOR method to make more it robust.
  for (let i = 0; i < DEFAULT_LENGTH; i++) {
    shortCode = BASE_62[Math.floor(temp % 62)] + shortCode;
    temp = Math.floor(temp / 62);
  }

  return shortCode;
};

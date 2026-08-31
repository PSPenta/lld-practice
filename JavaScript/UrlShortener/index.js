const express = require('express');
const { body, param, validationResult } = require('express-validator');
const app = express();
app.use(express.json());

const { generateShortCode } = require('./uniqueCode');
const DOMAIN = 'http://localhost:3001/';

const db = new Map();

app.post(
  '/shorten',
  [
    body('url')
      .trim()
      .notEmpty().withMessage('URL must not be empty!')
      .bail()
      .isURL({ protocols: [ 'http', 'https' ], require_protocol: true }).withMessage('URL must contain valid protocol!')
      .bail()
      .isLength({ max: 2000 }).withMessage('URL is too long!'),
    body('expiry')
      .optional({ nullable: true })
      .isInt({ min: 1, max: 365 }).withMessage('expiry must be between 1 and 365 days')
      .bail()
      .default(30)
  ],
  async (req, res, next) => {
    try {
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        return res.status(422).json({ success: false, errors: errors.array() });
      }

      const { url, expiry } = req.body; 

      const shortCode = generateShortCode(db);
      const shortUrl = DOMAIN + shortCode;

      db.set(shortCode, {
        longUrl: url,
        shortUrl,
        expiry: new Date(Date.now() + (expiry * 24 * 60 * 60 * 1000))
      });
      return res.status(201).json({ success: true, shortUrl });
    } catch (error) {
      return res.status(500).json({ success: false, error: 'Internal server error!' });
    }
  }
);

app.get(
  '/:shortCode',
  [
    param('shortCode')
      .notEmpty()
      .isAlphanumeric()
      .isLength({ min: 6, max: 6 })
  ],
  (req, res, next) => {
    try {
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        return res.status(422).send('<html><head><title>URL Shortener</title></head><body><h1>Invalid URL ;)</h1></body></html>');
      }

      const shortCode = req.params.shortCode;
      const uriDetails = db.get(shortCode);
      if (!uriDetails || !uriDetails.longUrl) {
        return res.status(400).send('<html><head><title>URL Shortener</title></head><body><h1>URL not found ;)</h1></body></html>');
      }

      if (uriDetails.expiry < Date.now()) {
        // If we use Redis or any other caching server, eviction will be done by it internally
        db.delete(shortCode);
        return res.status(400).send('<html><head><title>URL Shortener</title></head><body><h1>URL Expired ;)</h1></body></html>');
      }

      // We can change it to 301 status code if we want to get this redirection cached so that no hits will come to our backend for same URL.
      return res.status(302).redirect(uriDetails.longUrl);
    } catch (error) {
      return res.status(500).send('<html><head><title>URL Shortener</title></head><body><h1>Something went wrong, please try again later ;)</h1></body></html>');
    }
  }
);

app.listen(3001, () => console.log('Server started at port 3001'));

'use strict';

/** Rupees (float OK at API) → integer paise. All ledger math uses this. */
function toAmount(rupees) {
  if (typeof rupees !== 'number' || Number.isNaN(rupees)) {
    throw new Error('invalid amount');
  }
  return Math.round(rupees * 100);
}

function fromAmount(paise) {
  return paise / 100;
}

module.exports = { toAmount, fromAmount };

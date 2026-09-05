/**
 * @interface Expense (abstract base class)
 *
 * JavaScript has no native `interface` keyword. Shared expense logic lives here;
 * subclasses must implement {@link Expense#validate}.
 *
 * @method validate() → boolean
 *
 * @see ExactExpense, EqualExpense, PercentageExpense
 */
class Expense {
  constructor(id, paidBy, amount, splits) {
    if (!splits || !splits.length) {
      throw new Error('Expense must have at least one split');
    }

    this.id = id;
    this.paidBy = paidBy;
    this.amount = amount;
    this.splits = splits;

    return this;
  }

  validate() {
    throw new Error('Not implemented');
  }

  apply(BalanceSheet) {
    this.splits.forEach(split => {
      BalanceSheet.addDebt(split.user, this.paidBy, split.amount);
    });
  }
}

/** @implements {Expense} */
class ExactExpense extends Expense {
  validate() {
    const total = this.splits.reduce((sum, split) => sum + split.amount, 0);

    if (total !== this.amount) {
      throw new Error('Total expense does not match amount');
    }

    return this;
  }
}

/** @implements {Expense} */
class EqualExpense extends Expense {
  validate() {
    // amount is integer paise — floor + remainder on last (sum === total)
    const n = this.splits.length;
    // Create a base amount for each split without decimal
    const base = Math.floor(this.amount / n);
    // Assign the base amount to each split
    this.splits.forEach((split) => {
      split.amount = base;
    });
    // Add the remainder to the last split
    this.splits[n - 1].amount += this.amount - (base * n);

    return this;
  }
}

/** @implements {Expense} */
class PercentageExpense extends Expense {
  validate() {
    const totalPercentage = this.splits.reduce(
      (sum, split) => sum + split.percentage,
      0,
    );

    if (totalPercentage !== 100) {
      throw new Error('Total percentage must be 100');
    }

    // integer paise: floor each share, last gets remainder so sum === total
    let allocated = 0, smallest = this.amount, smallestIndex = 0;
    this.splits.forEach((split, i) => {
      // Assign the amount to the split based on the percentage
      split.amount = Math.floor((this.amount * split.percentage) / 100);
      // Add the amount to the allocated amount
      allocated += split.amount;
      if (smallest > split.amount) {
        smallest = split.amount;
        smallestIndex = i;
      }
    });

    // Add the remainder to the smallest split
    this.splits[smallestIndex].amount += this.amount - allocated;

    return this;
  }
}

class ExpenseFactory {
  static createExpense(expense) {
    const { id, type, paidBy, amount, splits } = expense;
    if (!type || !paidBy || !amount || !splits) {
      throw new Error('Invalid expense');
    }

    switch (type.toLowerCase()) {
      case 'exact':
        return new ExactExpense(id, paidBy, amount, splits);
      case 'equal':
        return new EqualExpense(id, paidBy, amount, splits);
      case 'percentage':
        return new PercentageExpense(id, paidBy, amount, splits);
      default:
        throw new Error('Invalid expense type');
    }
  }
}

module.exports = { ExpenseFactory };

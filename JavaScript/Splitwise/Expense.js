class Expense {
  constructor(paidBy, amount, splits) {
    if (!splits || !splits.length) {
      throw new Error('Expense must have at least one split');
    }

    this.paidBy = paidBy;
    this.amount = amount;
    this.splits = splits;
  }

  validate() {
    throw new Error('Not implemented');
  }

  apply(BalanceSheet) {
    this.splits.forEach((split) => {
      if (split.user !== this.paidBy) {
        BalanceSheet.addBalance(split.user, split.amount);
      } else {
        BalanceSheet.addBalance(split.user, -split.amount);
      }
    });
  }
}

class ExactExpense extends Expense {
  validate() {
    const total = this.splits.reduce((sum, split) => sum + split.amount, 0);

    if (total !== this.amount) {
      throw new Error('Total expense does not match amount');
    }

    return true;
  }
}

class EqualExpense extends Expense {
  validate() {
    const amount = Number((this.amount / this.splits.length).toFixed(2));

    this.splits.forEach((split) => {
      split.amount = amount;
    });

    return true;
  }
}

class PercentageExpense extends Expense {
  validate() {
    const totalPercentage = this.splits.reduce(
      (sum, split) => sum + split.percentage,
      0,
    );

    if (totalPercentage !== 100) {
      throw new Error('Total percentage must be 100');
    }

    this.splits.forEach((split) => {
      split.amount = Number(
        ((this.amount * split.percentage) / 100).toFixed(2),
      );
    });

    return true;
  }
}

class ExpenseFactory {
  static createExpense(type, paidBy, amount, splits) {
    switch (type.toLowerCase()) {
      case 'exact':
        return new ExactExpense(paidBy, amount, splits);
      case 'equal':
        return new EqualExpense(paidBy, amount, splits);
      case 'percentage':
        return new PercentageExpense(paidBy, amount, splits);
      default:
        throw new Error('Invalid expense type');
    }
  }
}

module.exports = { ExpenseFactory };

class Database {
  name = '';
  tables = new Map();

  constructor(name) {
    if (!name || typeof name !== 'string') throw new Error('Database name is invalid!');

    this.name = name;
  }

  showTables() {
    return Array.from(this.tables.keys());
  }

  addTable(table) {
    if (!table || !table.name || !table.properties || !table.properties.schema) {
      throw new Error('Invalid table name or schema!');
    }

    if (this.tables.has(table.name)) throw new Error(`Table ${table.name} already exists!`);

    this.tables.set(table.name, table);
  }

  getTable(name) {
    if (!name || typeof name !== 'string') throw new Error('Please enter valid table name!');

    if (!this.tables.has(name)) throw new Error(`Table ${name} doesn't exist!`);

    return this.tables.get(name);
  }

  deleteTable(name) {
    if (!name || typeof name !== 'string') throw new Error('Please enter valid table name!');

    if (!this.tables.has(name)) throw new Error(`Table ${name} doesn't exist!`);

    this.tables.delete(name);

    return true;
  }

  truncate() {
    this.tables.clear();
  }

  drop() {
    this.tables.clear();
    this.name = '';

    this = null;
  }
}

module.exports = { Database };

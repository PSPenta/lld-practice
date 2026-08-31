class Table {
  name = '';
  columnTypes = ['integer', 'string', 'boolean'];
  properties = {
    lastId: 0,
    schema: {}, // Schema contains columns and their types
    indexes: {}, // Indexes contains each column with their associated Maps
  };
  rows = new Map();

  constructor(tname, schema, indexes = null) {
    if (!tname || !schema || !schema.length) {
      throw new Error('Table name and schema cannot be empty!');
    }

    for (let { name, type, limit, required = false } of schema) {
      if (name in this.properties.schema) {
        throw new Error(`Column ${name} already exists!`);
      }

      if (!this.columnTypes.includes(type)) {
        throw new Error(`Unsupported column type ${type} for column ${name}!`);
      }

      this.properties.schema[name] = { type, limit, required };
    }
    this.name = tname;

    if (indexes) {
      if (!indexes || !indexes.length) throw new Error('Invalid indexes!');

      for (let col of indexes) {
        if (this.properties.indexes[col])
          throw new Error('Index already exists!');

        this.properties.indexes[col] = new Map();
      }
    }
  }

  addIndex(column) {
    if (!column || !this.properties.schema[column])
      throw new Error('Invalid column!');

    this.properties.indexes[column] = new Map();

    if (this.rows.size) {
      for (let [id, row] of this.rows) {
        if (!this.properties.indexes[column].has(row[column])) {
          this.properties.indexes[column].set(row[col], []);
        }
        this.properties.indexes[column].get(row[col]).push(id);
      }
    }
  }

  removeIndex(column) {
    if (!column || !this.properties.schema[column])
      throw new Error('Invalid column!');

    if (!this.properties.indexes[column]) return false;

    delete this.properties.indexes[column];

    return true;
  }

  select(where) {
    if (!where || !Object.keys(where)) return this.rows;

    let schema = Object.keys(this.properties.schema);
    let ids = [];
    let cols = Object.keys(where);
    for (let col of cols) {
      if (!schema.includes(col)) {
        throw new Error('Invalid column!');
      }

      if (this.properties.indexes[col]) {
        ids.push(...this.properties.indexes[col].get(where[col]));
        delete where[col];
      }
    }

    let result = [];
    ids = Array.from(new Set(ids));
    for (let id of ids) {
      result.push(this.rows.get(id));
    }

    if (Object.keys(where).length) {
      let rows = Array.from(this.rows.values());
      let result2 = rows.filter((row) => {
        let cols = Object.keys(where);
        for (const col in cols) {
          if (row[col] !== where[col]) {
            return false;
          }
        }
        return true;
      });
      result = [...result, ...result2];
    }

    return result;
  }

  insert(row) {
    if (!row || !Object.keys(row)) throw new Error('Invalid data to insert!');

    let schema = Object.keys(this.properties.schema);
    let rowSchema = Object.keys(row);
    for (let col of rowSchema) {
      if (!schema.includes(col)) {
        throw new Error('Invalid column!');
      }
    }

    let schemaKeys = Object.keys(this.properties.schema);
    for (let key of schemaKeys) {
      if (this.properties.schema[key].required && !row[key]) {
        throw new Error(`Column ${key} is a required field!`);
      }
    }

    let index = (this.properties.lastId || 0) + 1;
    this.rows.set(index, row);
    this.properties.lastId++;

    for (let col of rowSchema) {
      if (this.properties.indexes[col]) {
        if (!this.properties.indexes[col].has(row[col])) {
          this.properties.indexes[col].set(row[col], []);
        }
        this.properties.indexes[col].get(row[col]).push(index);
      }
    }
  }

  update(id, row) {
    if (!row || !Object.keys(row)) throw new Error('Invalid data to insert!');

    if (!id || typeof id !== 'number') throw new Error('Invalid row ID!');

    let schema = Object.keys(this.properties.schema);
    let rowSchema = Object.keys(row);
    for (let col of rowSchema) {
      if (!schema.includes(col)) {
        throw new Error('Invalid column!');
      }
    }

    this.rows.set(id, { ...this.rows.get(id), ...row });

    for (let col of rowSchema) {
      if (this.properties.indexes[col]) {
        if (!this.properties.indexes[col].has(row[col])) {
          this.properties.indexes[col].set(row[col], []);
        }
        if (!this.properties.indexes[col].get(row[col]).includes(id)) {
          this.properties.indexes[col].get(row[col]).push(id);
        }
      }
    }
  }

  delete(id) {
    if (!id || typeof id !== 'number') throw new Error('Invalid row!');

    if (!this.properties.schema[id]) {
      console.log("Row doesn't exist!");
      return;
    }

    this.properties.schema.delete(id);

    for (let [key, val] of this.properties.indexes) {
      if (val.includes(id)) delete val[val.indexOf(id)];
    }
  }

  truncate() {
    this.rows.clear();
    this.properties.lastId = 0;
  }

  drop() {
    this.rows.clear();
    this.properties.lastId = 0;
    this.properties.schema = {};
    this.properties.indexes = {};
    this.name = '';

    this = null;
  }
}

module.exports = { Table };

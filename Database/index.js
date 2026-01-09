const { Database } = require('./Database');
const { Table } = require('./Table');

try {
  let shoppingApp = new Database('shopping_app');

  shoppingApp.addTable(new Table('products', [
    { name: 'id', type: 'integer', limit: 10, required: false },
    { name: 'name', type: 'string', limit: 255, required: true },
    { name: 'price', type: 'integer', limit: 10, required: true },
    { name: 'quantity', type: 'integer', limit: 10, required: true },
  ]));

  shoppingApp.addTable(new Table('customers', [
    { name: 'id', type: 'integer', limit: 10, required: false },
    { name: 'name', type: 'string', limit: 255, required: true },
    { name: 'email', type: 'string', limit: 255, required: true },
    { name: 'address', type: 'string', limit: 255, required: true },
  ]));

  shoppingApp.addTable(new Table('orders', [
    { name: 'id', type: 'integer', limit: 10, required: false },
    { name: 'customer_id', type: 'integer', limit: 10, required: true },
    { name: 'total_price', type: 'integer', limit: 10, required: true },
    { name: 'order_date', type: 'string', limit: 255, required: true },
  ]));
  console.log('All Tables:', shoppingApp.showTables());

  let customers = shoppingApp.getTable('customers');
  customers.addIndex('email');

  customers.insert({
    name: 'John Doe',
    email: 'john.doe@example.com',
    address: '123 Main St',
  });
  console.log('Customer Rows:', customers.rows);
  console.log('Customer Indexes:', customers.indexes);

  customers.update(1, { address: '456 Elm St' });
  console.log(customers.rows);
  customers.insert({
    name: 'Pritesh Shinde',
    email: 'shindepritesh78@gmail.com',
    address: 'Navi Mumbai',
  });
  console.log('Customer Indexes:', customers.indexes);

  customers.insert({
    name: 'Pritesh Shinde',
    email: 'pritesh.shinde@medibuddy.in',
    address: 'Mumbai',
  });
  console.log('Customer Rows:', customers.rows);
  console.log('Customer Indexes:', customers.indexes);

  let res = customers.select({});
  console.log('Result 1: ', res);

  res = customers.select({ name: 'Pritesh Shinde' });
  console.log('Result 2: ', res);

  res = customers.select({ email: 'pritesh.shinde@medibuddy.in' });
  console.log('Result 3: ', res);

  res = customers.select({ email: 'pritesh.shinde@medibuddy.in', name: 'Pritesh Shind' });
  console.log('Result 4: ', res);

  customers.delete(1);
  console.log('Customer Rows:', customers.rows);
  console.log('Customer Indexes:', customers.indexes);

  let products = shoppingApp.getTable('products');
  products.insert({
    name: 'Product 1',
    price: 10,
    quantity: 100,
  });
  console.log('Product Rows:', products.rows);

  let orders = shoppingApp.getTable('orders');
  orders.insert({
    customer_id: 1,
    total_price: 100,
    order_date: '2022-01-01',
  });
  console.log('Order Rows:', orders.rows);
} catch (error) {
  console.error(error);
}

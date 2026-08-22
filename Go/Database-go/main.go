package main

import "fmt"

func main() {
	shoppingApp, err := NewDatabase("shopping_app")
	if err != nil {
		fmt.Println(err)
		return
	}

	products, _ := NewTable("products", []ColumnDef{
		{Name: "id", Type: ColumnInteger, Limit: 10},
		{Name: "name", Type: ColumnString, Limit: 255, Required: true},
		{Name: "price", Type: ColumnInteger, Limit: 10, Required: true},
		{Name: "quantity", Type: ColumnInteger, Limit: 10, Required: true},
	}, nil)

	customers, _ := NewTable("customers", []ColumnDef{
		{Name: "id", Type: ColumnInteger, Limit: 10},
		{Name: "name", Type: ColumnString, Limit: 255, Required: true},
		{Name: "email", Type: ColumnString, Limit: 255, Required: true},
		{Name: "address", Type: ColumnString, Limit: 255, Required: true},
	}, nil)

	orders, _ := NewTable("orders", []ColumnDef{
		{Name: "id", Type: ColumnInteger, Limit: 10},
		{Name: "customer_id", Type: ColumnInteger, Limit: 10, Required: true},
		{Name: "total_price", Type: ColumnInteger, Limit: 10, Required: true},
		{Name: "order_date", Type: ColumnString, Limit: 255, Required: true},
	}, nil)

	shoppingApp.AddTable(products)
	shoppingApp.AddTable(customers)
	shoppingApp.AddTable(orders)
	fmt.Println("All Tables:", shoppingApp.ShowTables())

	customersTable, _ := shoppingApp.GetTable("customers")
	customersTable.AddIndex("email")

	customersTable.Insert(map[string]any{
		"name":    "John Doe",
		"email":   "john.doe@example.com",
		"address": "123 Main St",
	})
	fmt.Println("Customer Rows:", customersTable.Rows)

	customersTable.Update(1, map[string]any{"address": "456 Elm St"})
	fmt.Println(customersTable.Rows)

	customersTable.Insert(map[string]any{
		"name":    "Pritesh Shinde",
		"email":   "shindepritesh78@gmail.com",
		"address": "Navi Mumbai",
	})

	customersTable.Insert(map[string]any{
		"name":    "Pritesh Shinde",
		"email":   "pritesh.shinde@medibuddy.in",
		"address": "Mumbai",
	})
	fmt.Println("Customer Rows:", customersTable.Rows)

	res, _ := customersTable.Select(map[string]any{})
	fmt.Println("Result 1: ", res)

	res, _ = customersTable.Select(map[string]any{"name": "Pritesh Shinde"})
	fmt.Println("Result 2: ", res)

	res, _ = customersTable.Select(map[string]any{"email": "pritesh.shinde@medibuddy.in"})
	fmt.Println("Result 3: ", res)

	res, _ = customersTable.Select(map[string]any{"email": "pritesh.shinde@medibuddy.in", "name": "Pritesh Shind"})
	fmt.Println("Result 4: ", res)

	customersTable.Delete(1)
	fmt.Println("Customer Rows:", customersTable.Rows)

	productsTable, _ := shoppingApp.GetTable("products")
	productsTable.Insert(map[string]any{
		"name":     "Product 1",
		"price":    10,
		"quantity": 100,
	})
	fmt.Println("Product Rows:", productsTable.Rows)

	ordersTable, _ := shoppingApp.GetTable("orders")
	ordersTable.Insert(map[string]any{
		"customer_id": 1,
		"total_price": 100,
		"order_date":  "2022-01-01",
	})
	fmt.Println("Order Rows:", ordersTable.Rows)
}

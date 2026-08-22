package main

import "fmt"

type Database struct {
	Name   string
	Tables map[string]*Table
}

func NewDatabase(name string) (*Database, error) {
	if name == "" {
		return nil, fmt.Errorf("database name is invalid!")
	}
	return &Database{Name: name, Tables: make(map[string]*Table)}, nil
}

func (d *Database) ShowTables() []string {
	names := make([]string, 0, len(d.Tables))
	for name := range d.Tables {
		names = append(names, name)
	}
	return names
}

func (d *Database) AddTable(table *Table) error {
	if table == nil || table.Name == "" || len(table.Properties.Schema) == 0 {
		return fmt.Errorf("invalid table name or schema!")
	}
	if _, exists := d.Tables[table.Name]; exists {
		return fmt.Errorf("table %s already exists!", table.Name)
	}
	d.Tables[table.Name] = table
	return nil
}

func (d *Database) GetTable(name string) (*Table, error) {
	if name == "" {
		return nil, fmt.Errorf("please enter valid table name!")
	}
	table, ok := d.Tables[name]
	if !ok {
		return nil, fmt.Errorf("table %s doesn't exist!", name)
	}
	return table, nil
}

func (d *Database) DeleteTable(name string) (bool, error) {
	if name == "" {
		return false, fmt.Errorf("please enter valid table name!")
	}
	if _, ok := d.Tables[name]; !ok {
		return false, fmt.Errorf("table %s doesn't exist!", name)
	}
	delete(d.Tables, name)
	return true, nil
}

func (d *Database) Truncate() {
	d.Tables = make(map[string]*Table)
}

func (d *Database) Drop() {
	d.Truncate()
	d.Name = ""
}

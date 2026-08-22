package main

import "fmt"

type ColumnType string

const (
	ColumnInteger ColumnType = "integer"
	ColumnString  ColumnType = "string"
	ColumnBoolean ColumnType = "boolean"
)

type ColumnDef struct {
	Name     string
	Type     ColumnType
	Limit    int
	Required bool
}

type ColumnSchema struct {
	Type     ColumnType
	Limit    int
	Required bool
}

type TableProperties struct {
	LastID  int
	Schema  map[string]ColumnSchema
	Indexes map[string]map[any][]int
}

type Table struct {
	Name       string
	Properties TableProperties
	Rows       map[int]map[string]any
}

func NewTable(name string, schema []ColumnDef, indexes []string) (*Table, error) {
	if name == "" || len(schema) == 0 {
		return nil, fmt.Errorf("table name and schema cannot be empty!")
	}

	columnTypes := map[ColumnType]bool{
		ColumnInteger: true,
		ColumnString:  true,
		ColumnBoolean: true,
	}

	t := &Table{
		Name: name,
		Properties: TableProperties{
			Schema:  make(map[string]ColumnSchema),
			Indexes: make(map[string]map[any][]int),
		},
		Rows: make(map[int]map[string]any),
	}

	for _, col := range schema {
		if _, exists := t.Properties.Schema[col.Name]; exists {
			return nil, fmt.Errorf("column %s already exists!", col.Name)
		}
		if !columnTypes[col.Type] {
			return nil, fmt.Errorf("unsupported column type %s for column %s!", col.Type, col.Name)
		}
		t.Properties.Schema[col.Name] = ColumnSchema{
			Type:     col.Type,
			Limit:    col.Limit,
			Required: col.Required,
		}
	}

	for _, col := range indexes {
			if _, exists := t.Properties.Indexes[col]; exists {
				return nil, fmt.Errorf("index already exists!")
			}
			t.Properties.Indexes[col] = make(map[any][]int)
	}

	return t, nil
}

func (t *Table) AddIndex(column string) error {
	if _, ok := t.Properties.Schema[column]; !ok {
		return fmt.Errorf("invalid column!")
	}
	t.Properties.Indexes[column] = make(map[any][]int)
	for id, row := range t.Rows {
		val := row[column]
		t.Properties.Indexes[column][val] = append(t.Properties.Indexes[column][val], id)
	}
	return nil
}

func (t *Table) RemoveIndex(column string) (bool, error) {
	if _, ok := t.Properties.Schema[column]; !ok {
		return false, fmt.Errorf("invalid column!")
	}
	if _, ok := t.Properties.Indexes[column]; !ok {
		return false, nil
	}
	delete(t.Properties.Indexes, column)
	return true, nil
}

func (t *Table) Select(where map[string]any) ([]map[string]any, error) {
	if len(where) == 0 {
		result := make([]map[string]any, 0, len(t.Rows))
		for _, row := range t.Rows {
			result = append(result, row)
		}
		return result, nil
	}

	remaining := make(map[string]any)
	for k, v := range where {
		remaining[k] = v
	}

	var ids []int
	for col, val := range remaining {
		if _, ok := t.Properties.Schema[col]; !ok {
			return nil, fmt.Errorf("invalid column!")
		}
		if index, ok := t.Properties.Indexes[col]; ok {
			ids = append(ids, index[val]...)
			delete(remaining, col)
		}
	}

	idSet := make(map[int]bool)
	for _, id := range ids {
		idSet[id] = true
	}

	var result []map[string]any
	for id := range idSet {
		if row, ok := t.Rows[id]; ok {
			result = append(result, row)
		}
	}

	if len(remaining) > 0 {
		for _, row := range t.Rows {
			match := true
			for col, val := range remaining {
				if row[col] != val {
					match = false
					break
				}
			}
			if match {
				result = append(result, row)
			}
		}
	}

	return result, nil
}

func (t *Table) Insert(row map[string]any) error {
	if len(row) == 0 {
		return fmt.Errorf("invalid data to insert!")
	}

	for col := range row {
		if _, ok := t.Properties.Schema[col]; !ok {
			return fmt.Errorf("invalid column!")
		}
	}

	for key, schema := range t.Properties.Schema {
		if schema.Required {
			if val, ok := row[key]; !ok || val == nil || val == "" {
				return fmt.Errorf("column %s is a required field!", key)
			}
		}
	}

	id := t.Properties.LastID + 1
	t.Rows[id] = row
	t.Properties.LastID = id

	for col := range row {
		if index, ok := t.Properties.Indexes[col]; ok {
			val := row[col]
			index[val] = append(index[val], id)
		}
	}

	return nil
}

func (t *Table) Update(id int, row map[string]any) error {
	if len(row) == 0 {
		return fmt.Errorf("invalid data to update!")
	}
	if id <= 0 {
		return fmt.Errorf("invalid row ID!")
	}

	for col := range row {
		if _, ok := t.Properties.Schema[col]; !ok {
			return fmt.Errorf("invalid column!")
		}
	}

	existing := t.Rows[id]
	merged := make(map[string]any)
	for k, v := range existing {
		merged[k] = v
	}
	for k, v := range row {
		merged[k] = v
	}
	t.Rows[id] = merged

	for col := range row {
		if index, ok := t.Properties.Indexes[col]; ok {
			val := row[col]
			found := false
			for _, existingID := range index[val] {
				if existingID == id {
					found = true
					break
				}
			}
			if !found {
				index[val] = append(index[val], id)
			}
		}
	}

	return nil
}

func (t *Table) Delete(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid row!")
	}
	if _, ok := t.Rows[id]; !ok {
		return fmt.Errorf("row %d does not exist", id)
	}
	delete(t.Rows, id)
	return nil
}

func (t *Table) Truncate() {
	t.Rows = make(map[int]map[string]any)
	t.Properties.LastID = 0
}

func (t *Table) Drop() {
	t.Truncate()
	t.Properties.Schema = make(map[string]ColumnSchema)
	t.Properties.Indexes = make(map[string]map[any][]int)
	t.Name = ""
}

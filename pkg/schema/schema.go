package schema

// Schema represents a database schema.
type Schema struct {
	Name   string
	Tables []*Table
}

// Table represents a database table.
type Table struct {
	Name    string
	Columns []*Column
}

// Column represents a database column.
type Column struct {
	Name string
	Type string
}

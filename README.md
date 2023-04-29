# Relational Database Schemas Module

This Go module provides a set of functions and types for working with relational database schemas, tables, and column types. It also includes implementations for popular database drivers to allow users to easily retrieve schema information for their databases.

## Usage

```
go get github.com/ItsAminZamani/dbmeta@0.1.0
```

### Retrieving a schema

To retrieve a schema for a database, use the `GetSchema` function of the appropriate driver:

```go
import (
    "github.com/ItsAminZamani/dbmeta"
)

func main() {
    driver, err := dbmeta.NewDriver("mysql", "localhost", 3306, "root", "my-secret-pw", "example_db")
	if err != nil {
		// Handle error
	}
    schemas, err := dbmeta.GetSchemas(driver)
	if err != nil {
        // Handle error
    }

    // Use schema object
    for _, schema := range schemas {
		fmt.Println("Schema:", schema.Name)
		for _, table := range schema.Tables {
			fmt.Println("  Table:", table.Name)
			for _, column := range table.Columns {
				fmt.Printf("    Column: %s (%s)\n", column.Name, column.Type)
			}
		}
	}
}

```

### Retrieving tables from a schema

To retrieve tables from a schema, use the `GetTables` function of the appropriate driver:

```go
import (
    "github.com/ItsAminZamani/dbmeta"
)

func main() {
    driver, err := dbmeta.NewDriver("postgres", "localhost", 5432, "root", "my-secret-pw", "example_db")
	if err != nil {
		// Handle error
	}
    tables, err := dbmeta.GetTables(driver, "public")
    if err != nil {
        // Handle error
    }
}
```

### Retrieving columns from a table

To retrieve columns from a table, use the `GetColumns` function of the appropriate driver:

```go
import (
    "github.com/ItsAminZamani/dbmeta"
)

func main() {
    driver, err := dbmeta.NewDriver("mssql", "localhost", 1433, "root", "my-secret-pw", "example_db")
    if err != nil {
        // Handle error
    }
    columns, err := dbmeta.GetColumns(driver, "dbo", "users")
    if err != nil {
        // Handle error
    }
}
```

### Supported database drivers

Currently, the following database drivers are supported:

- MySQL (github.com/ItsAminZamani/dbmeta/pkg/drivers/mysql)
- PostgreSQL (github.com/ItsAminZamani/dbmeta/pkg/drivers/postgres)
- SQL Server (github.com/ItsAminZamani/dbmeta/pkg/drivers/mssql)

More database drivers may be added in the future.

## Schema structure

The schema package provides the following types for representing database schemas:

```go
type Schema struct {
    Name   string
    Tables []*Table
}

type Table struct {
    Name    string
    Columns []*Column
}

type Column struct {
    Name string
    Type string
}
```

**Schema** represents an entire database schema, and contains a list of Table objects.

**Table** represents a single database table, and contains a list of Column objects.

**Column** represents a single column in a database table, and contains the name of the column and its data type.

## Contributing

Contributions to this module are welcome! If you'd like to add support for a new database driver or make any other changes, please fork this repository and submit a pull request.

## License

This module is licensed under the MIT License. See LICENSE file for details.

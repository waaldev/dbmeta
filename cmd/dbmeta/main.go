package main

import (
	"fmt"

	"github.com/ItsAminZamani/dbmeta"
)

func main() {
	driver, err := dbmeta.NewDriver("mysql", "localhost", 3306, "root", "my-secret-pw", "example_db")
	if err != nil {
		panic(err)
	}
	schemas, err := dbmeta.GetSchemas(driver)
	if err != nil {
		panic(err)
	}
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

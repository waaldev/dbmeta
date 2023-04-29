package mssql

import (
	"database/sql"
	"fmt"

	"github.com/ItsAminZamani/dbmeta/pkg/schema"
	_ "github.com/denisenkom/go-mssqldb"
)

// SQLServerDriver is a driver for SQL Server databases.
type MSSQLDriver struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
}

// GetSchemas returns the schemas for a SQL Server database.
func (d *MSSQLDriver) GetSchemas() ([]*schema.Schema, error) {
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", d.Host, d.Username, d.Password, d.Port, d.DatabaseName))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT name FROM sys.schemas WHERE name NOT IN ('guest', 'INFORMATION_SCHEMA', 'sys')")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []*schema.Schema
	for rows.Next() {
		var schemaName string
		if err := rows.Scan(&schemaName); err != nil {
			return nil, err
		}

		schemas = append(schemas, &schema.Schema{Name: schemaName})
	}

	for _, schema := range schemas {
		if err := d.getTables(db, schema); err != nil {
			return nil, err
		}
	}

	return schemas, nil
}

// getTables returns the tables for a SQL Server schema.
func (d *MSSQLDriver) getTables(db *sql.DB, sch *schema.Schema) error {
	rows, err := db.Query("SELECT name FROM sys.tables WHERE schema_id = SCHEMA_ID(?)", sch.Name)
	if err != nil {
		return err
	}
	defer rows.Close()
	var tables []*schema.Table
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return err
		}

		tables = append(tables, &schema.Table{Name: tableName})
	}

	for _, table := range tables {
		if err := d.getColumns(db, sch, table); err != nil {
			return err
		}
	}

	sch.Tables = tables

	return nil
}

// getColumns returns the columns for a SQL Server table.
func (d *MSSQLDriver) getColumns(db *sql.DB, sch *schema.Schema, tbl *schema.Table) error {
	rows, err := db.Query("SELECT name, type_name(user_type_id) FROM sys.columns WHERE object_id = OBJECT_ID(?)", sch.Name+"."+tbl.Name)
	if err != nil {
		return err
	}
	defer rows.Close()
	var columns []*schema.Column
	for rows.Next() {
		var columnName string
		var typeName string
		if err := rows.Scan(&columnName, &typeName); err != nil {
			return err
		}

		columns = append(columns, &schema.Column{Name: columnName, Type: typeName})
	}

	tbl.Columns = columns

	return nil
}

// GetTableNames returns the table names for a SQL Server schema.
func (d *MSSQLDriver) GetTables(schemaName string) ([]*schema.Table, error) {
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", d.Host, d.Username, d.Password, d.Port, d.DatabaseName))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT name FROM sys.tables WHERE schema_id = SCHEMA_ID(?)", schemaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []*schema.Table
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, &schema.Table{Name: tableName})
	}
	return tables, nil
}

// GetColumnNames returns the column names for a SQL Server table.
func (d *MSSQLDriver) GetColumns(schemaName, tableName string) ([]*schema.Column, error) {
	db, err := sql.Open("mssql", fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", d.Host, d.Username, d.Password, d.Port, d.DatabaseName))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT name, type_name(user_type_id) FROM sys.columns WHERE object_id = OBJECT_ID(?)", schemaName+"."+tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var columns []*schema.Column
	for rows.Next() {
		var columnName string
		var typeName string
		if err := rows.Scan(&columnName, &typeName); err != nil {
			return nil, err
		}
		columns = append(columns, &schema.Column{Name: columnName, Type: typeName})
	}
	return columns, nil
}

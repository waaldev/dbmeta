package mysql

import (
	"database/sql"
	"fmt"

	"github.com/ItsAminZamani/dbmeta/pkg/schema"
	_ "github.com/go-sql-driver/mysql"
)

// MySQLDriver is a driver for MySQL databases.
type MySQLDriver struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
}

// GetSchemas returns the schemas for a MySQL database.
func (d *MySQLDriver) GetSchemas() ([]*schema.Schema, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.Username, d.Password, d.Host, d.Port, d.DatabaseName))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')")
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

// getTables returns the tables for a MySQL schema.
func (d *MySQLDriver) getTables(db *sql.DB, sch *schema.Schema) error {
	rows, err := db.Query("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ?", sch.Name)
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

// getColumns returns the columns for a MySQL table.
func (d *MySQLDriver) getColumns(db *sql.DB, sch *schema.Schema, tbl *schema.Table) error {
	rows, err := db.Query("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?", sch.Name, tbl.Name)
	if err != nil {
		return err
	}
	defer rows.Close()
	var columns []*schema.Column
	for rows.Next() {
		var columnName, columnType string
		if err := rows.Scan(&columnName, &columnType); err != nil {
			return err
		}

		columns = append(columns, &schema.Column{Name: columnName, Type: columnType})
	}
	tbl.Columns = columns
	return nil
}

// GetTables returns the tables for a MySQL schema.
func (d *MySQLDriver) GetTables(schemaName string) ([]*schema.Table, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.Username, d.Password, d.Host, d.Port, d.DatabaseName))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ?", schemaName)
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

// GetColumns returns the columns for a MySQL table.
func (d *MySQLDriver) GetColumns(schemaName, tableName string) ([]*schema.Column, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.Username, d.Password, d.Host, d.Port, d.DatabaseName))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT COLUMN_NAME, DATA_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?", schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var columns []*schema.Column
	for rows.Next() {
		var columnName, columnType string
		if err := rows.Scan(&columnName, &columnType); err != nil {
			return nil, err
		}

		columns = append(columns, &schema.Column{Name: columnName, Type: columnType})
	}
	return columns, nil
}

// TestConnection tests the connection to a Postgres database.
func (d *MySQLDriver) TestConnection() error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.Username, d.Password, d.Host, d.Port, d.DatabaseName))
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

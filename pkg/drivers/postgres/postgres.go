package postgres

import (
	"database/sql"
	"fmt"

	"github.com/ItsAminZamani/dbmeta/pkg/schema"
	_ "github.com/lib/pq"
)

// PostgresDriver is a driver for Postgres databases.
type PostgresDriver struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
}

// GetSchemas returns the schemas for a Postgres database.
func (d *PostgresDriver) GetSchemas() ([]*schema.Schema, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.Username,
		d.Password, d.Host, d.Port, d.DatabaseName))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('information_schema', 'pg_catalog', 'pg_toast')")
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

// getTables returns the tables for a Postgres schema.
func (d *PostgresDriver) getTables(db *sql.DB, sch *schema.Schema) error {
	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = $1", sch.Name)
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

// getColumns returns the columns for a Postgres table.
func (d *PostgresDriver) getColumns(db *sql.DB, sch *schema.Schema, tbl *schema.Table) error {
	rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2", sch.Name, tbl.Name)
	if err != nil {
		return err
	}
	defer rows.Close()
	var columns []*schema.Column
	for rows.Next() {
		var columnName, dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			return err
		}

		columns = append(columns, &schema.Column{Name: columnName, Type: dataType})
	}

	tbl.Columns = columns

	return nil
}

// GetTables returns the tables for a MySQL schema.
func (d *PostgresDriver) GetTables(schemaName string) ([]*schema.Table, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.Username,
		d.Password, d.Host, d.Port, d.DatabaseName))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = $1", schemaName)
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

	for _, table := range tables {
		if err := d.getColumns(db, &schema.Schema{Name: schemaName}, table); err != nil {
			return nil, err
		}
	}

	return tables, nil
}

// GetColumns returns the columns for a MySQL table.
func (d *PostgresDriver) GetColumns(schemaName, tableName string) ([]*schema.Column, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.Username,
		d.Password, d.Host, d.Port, d.DatabaseName))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_schema = $1 AND table_name = $2", schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var columns []*schema.Column
	for rows.Next() {
		var columnName, dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			return nil, err
		}

		columns = append(columns, &schema.Column{Name: columnName, Type: dataType})
	}

	return columns, nil
}

// TestConnection tests the connection to a Postgres database.
func (d *PostgresDriver) TestConnection() error {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.Username,
		d.Password, d.Host, d.Port, d.DatabaseName))
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}
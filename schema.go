package dbmeta

import (
	"github.com/siaminz/dbmeta/pkg/drivers/mssql"
	"github.com/siaminz/dbmeta/pkg/drivers/mysql"
	"github.com/siaminz/dbmeta/pkg/drivers/postgres"
	"github.com/siaminz/dbmeta/pkg/schema"
)

type ErrInvalidDriver struct {
	Driver string
}

func (e *ErrInvalidDriver) Error() string {
	return "invalid driver: " + e.Driver
}

// Create a new driver for a database.
func NewDriver(driver string, host string, port int, username string, password string, databaseName string) (interface{}, error) {
	switch driver {
	case "mysql":
		return &mysql.MySQLDriver{
			Host:         host,
			Port:         port,
			Username:     username,
			Password:     password,
			DatabaseName: databaseName,
		}, nil
	case "postgres":
		return &postgres.PostgresDriver{
			Host:         host,
			Port:         port,
			Username:     username,
			Password:     password,
			DatabaseName: databaseName,
		}, nil
	case "mssql":
		return &mssql.MSSQLDriver{
			Host:         host,
			Port:         port,
			Username:     username,
			Password:     password,
			DatabaseName: databaseName,
		}, nil
	default:
		return nil, &ErrInvalidDriver{Driver: driver}
	}
}

// Get the schemas for a database.
func GetSchemas(driver interface{}) ([]*schema.Schema, error) {
	switch d := driver.(type) {
	case *mysql.MySQLDriver:
		return d.GetSchemas()
	case *postgres.PostgresDriver:
		return d.GetSchemas()
	case *mssql.MSSQLDriver:
		return d.GetSchemas()
	default:
		return nil, &ErrInvalidDriver{Driver: "unknown"}
	}
}

// Get the tables for a schema.
func GetTables(driver interface{}, schemaName string) ([]*schema.Table, error) {
	switch d := driver.(type) {
	case *mysql.MySQLDriver:
		return d.GetTables(schemaName)
	case *postgres.PostgresDriver:
		return d.GetTables(schemaName)
	case *mssql.MSSQLDriver:
		return d.GetTables(schemaName)
	default:
		return nil, &ErrInvalidDriver{Driver: "unknown"}
	}
}

// Get the columns for a table.
func GetColumns(driver interface{}, schemaName string, tableName string) ([]*schema.Column, error) {
	switch d := driver.(type) {
	case *mysql.MySQLDriver:
		return d.GetColumns(schemaName, tableName)
	case *postgres.PostgresDriver:
		return d.GetColumns(schemaName, tableName)
	case *mssql.MSSQLDriver:
		return d.GetColumns(schemaName, tableName)
	default:
		return nil, &ErrInvalidDriver{Driver: "unknown"}
	}
}

func TestConnection(driver interface{}) error {
	switch d := driver.(type) {
	case *mysql.MySQLDriver:
		return d.TestConnection()
	case *postgres.PostgresDriver:
		return d.TestConnection()
	case *mssql.MSSQLDriver:
		return d.TestConnection()
	default:
		return &ErrInvalidDriver{Driver: "unknown"}
	}
}

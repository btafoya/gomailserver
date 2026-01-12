package database

type Driver string

const (
	DriverSQLite   Driver = "sqlite3"
	DriverPostgres Driver = "postgres"
)

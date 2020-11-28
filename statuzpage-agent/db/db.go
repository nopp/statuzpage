package db

import (
	"database/sql"
	"statuzpage-agent/configuration"
)

// DBConnection responsible to return db connection
func DBConnection() (*sql.DB, error) {

	config := configuration.LoadConfiguration()
	db, err := sql.Open("mysql", ""+config.MySQLUser+":"+config.MySQLPass+"@tcp("+config.MySQLHost+")/"+config.MySQLDb)

	return db, err
}

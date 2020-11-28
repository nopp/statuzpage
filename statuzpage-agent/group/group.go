package group

import (
	"log"
	"statuzpage-agent/db"
)

// Group struct
type Group struct {
	ID   int
	Name string
}

// ReturnGroupName Return Group name from ID Group
func ReturnGroupName(IDGroup int) string {

	var name string

	db, errDB := db.DBConnection()
	defer db.Close()
	if errDB != nil {
		log.Println("Cant connect to server host!")
	}

	err := db.QueryRow("SELECT name FROM sp_groups WHERE id = ?", IDGroup).Scan(&name)

	if err != nil {
		log.Printf(err.Error())
	}

	return name
}

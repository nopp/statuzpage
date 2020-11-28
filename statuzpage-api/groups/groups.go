package groups

import (
	"encoding/json"
	"log"
	"net/http"
	"statuzpage-api/common"
	"strconv"

	"github.com/gorilla/mux"
)

type Group struct {
	ID                 string
	Name               string
	TotalIncidentsOpen string
}

// Group
func GetGroups(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		var groups []Group
		var group Group

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		rows, err := db.Query("SELECT g.id,g.name, (SELECT count(*) FROM sp_incidents i WHERE i.idGroup = g.id AND i.finishedat IS NULL) as open FROM sp_groups g")
		if err != nil {
			common.Message(w, "Cant get groups from systems!")
		}

		for rows.Next() {
			err := rows.Scan(&group.ID, &group.Name, &group.TotalIncidentsOpen)
			if err != nil {
				common.Message(w, "Cant return group informations!")
			}
			groups = append(groups, group)
		}

		json.NewEncoder(w).Encode(groups)
	} else {
		common.Message(w, "Invalid token!")
	}
}

func GetGroup(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		params := mux.Vars(r)
		var group Group

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		err := db.QueryRow("SELECT id,name from sp_groups WHERE id = ?", params["id"]).Scan(&group.ID, &group.Name)
		if err != nil {
			common.Message(w, "Cant get group from systems!")
		}

		json.NewEncoder(w).Encode(group)
	} else {
		common.Message(w, "Invalid token!")
	}
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		var group Group
		var total int

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		_ = json.NewDecoder(r.Body).Decode(&group)

		err := db.QueryRow("SELECT COUNT(*) from sp_groups WHERE name = ?", group.Name).Scan(&total)
		if err != nil {
			common.Message(w, "Cant count groups!")
		}

		if total == 0 {
			stmt, err := db.Prepare("INSERT INTO sp_groups(name) values(?)")
			if err != nil {
				common.Message(w, "Cant prepare insert group!")
			}
			res, err := stmt.Exec(group.Name)
			if err != nil {
				common.Message(w, "Cant insert group!")
			}

			lastID, _ := res.LastInsertId()
			group.ID = strconv.FormatInt(lastID, 10)

			json.NewEncoder(w).Encode(group)
		} else {
			common.Message(w, "Group already in database!")
		}
	} else {
		common.Message(w, "Invalid token!")
	}
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		params := mux.Vars(r)

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		stmt, err := db.Prepare("DELETE FROM sp_groups WHERE id = ?")
		if err != nil {
			common.Message(w, "Cant prepare delete group!")
		}

		_, err = stmt.Exec(params["id"])
		if err != nil {
			common.Message(w, "Cant delete group!")
		}

		common.Message(w, "Group "+params["id"]+" deleted!")
	} else {
		common.Message(w, "Invalid token!")
	}
}

func ReturnGroupName(IDGroup int) string {

	var name string

	db, errDB := common.DBConnection()
	defer db.Close()
	if errDB != nil {
		log.Println("Cant connect to server host!")
	}

	err := db.QueryRow("SELECT name FROM sp_groups WHERE id = ?", IDGroup).Scan(&name)

	if err != nil {
		log.Printf("Cant get group info!")
	}

	return name
}

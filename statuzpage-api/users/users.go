package users

import (
	"encoding/json"
	"log"
	"net/http"
	"statuzpage-api/common"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID       string
	Name     string
	Login    string
	Password string
}

// Users
func GetUsers(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		var users []User
		var user User

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {

			log.Println("Cant connect to server host!")
		}

		rows, err := db.Query("SELECT id,name,login,password from sp_users")
		if err != nil {
			common.Message(w, "Cant get users from systems!")
		}

		for rows.Next() {
			err := rows.Scan(&user.ID, &user.Name, &user.Login, &user.Password)
			if err != nil {
				common.Message(w, "Cant return user informations!")
			}
			users = append(users, user)
		}

		json.NewEncoder(w).Encode(users)
	} else {
		common.Message(w, "Invalid token!")
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		params := mux.Vars(r)
		var user User

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		err := db.QueryRow("SELECT id,name,login,password from sp_users WHERE id = ?", params["id"]).Scan(&user.ID, &user.Name, &user.Login, &user.Password)
		if err != nil {
			common.Message(w, "Cant get user from systems!")
		}

		json.NewEncoder(w).Encode(user)
	} else {
		common.Message(w, "Invalid token!")
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		var user User
		var total int

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			log.Println("Cant connect to server host!")
		}

		_ = json.NewDecoder(r.Body).Decode(&user)

		err := db.QueryRow("SELECT COUNT(*) from sp_users WHERE login = ?", user.Login).Scan(&total)
		if err != nil {
			common.Message(w, "Cant get user from systems!")
		}

		if total == 0 {
			stmt, err := db.Prepare("INSERT INTO sp_users(name,login,password) values(?,?,?)")
			if err != nil {
				common.Message(w, "Cant prepare insert user!")
			}

			res, err := stmt.Exec(user.Name, user.Login, user.Password)
			if err != nil {
				common.Message(w, "Cant insert user!")
			}

			lastID, _ := res.LastInsertId()
			user.ID = strconv.FormatInt(lastID, 10)

			json.NewEncoder(w).Encode(user)
		} else {
			common.Message(w, "User already in database!")
		}
	} else {
		common.Message(w, "Invalid token!")
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		params := mux.Vars(r)

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		stmt, err := db.Prepare("DELETE FROM sp_users WHERE id = ?")
		if err != nil {
			common.Message(w, "Cant prepare delete user!")
		}

		_, err = stmt.Exec(params["id"])
		if err != nil {
			common.Message(w, "Cant delete user!")
		}

		common.Message(w, "User "+params["id"]+" deleted!")
	} else {
		common.Message(w, "Invalid token!")
	}
}

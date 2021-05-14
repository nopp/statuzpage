package urls

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"statuzpage-api/common"

	"github.com/gorilla/mux"
)

type Url struct {
	ID            int
	Name          string
	URL           string
	ReturnCode    string
	Content       sql.NullString
	CheckInterval uint64
}

// Url
func GetUrls(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		var urls []Url
		var url Url

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		rows, err := db.Query("SELECT id,name,url,return_code,content,check_interval from sp_urls")
		if err != nil {
			common.Message(w, "Cant get urls from systems!")
		}

		for rows.Next() {
			err := rows.Scan(&url.ID, &url.Name, &url.URL, &url.ReturnCode, &url.Content, &url.CheckInterval)
			if err != nil {
				common.Message(w, "Cant return url informations!")
			}
			urls = append(urls, url)
		}

		json.NewEncoder(w).Encode(urls)
	} else {
		common.Message(w, "Invalid token!")
	}
}

func GetUrl(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		params := mux.Vars(r)
		var url Url

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		err := db.QueryRow("SELECT id,name,url,return_code,content,check_interval from sp_urls WHERE id = ?", params["id"]).Scan(&url.ID, &url.Name, &url.URL, &url.ReturnCode, &url.Content, &url.CheckInterval)
		if err != nil {
			common.Message(w, "Cant get url from systems!")
		}

		json.NewEncoder(w).Encode(url)
	} else {
		common.Message(w, "Invalid token!")
	}
}

func CreateUrl(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		var url Url
		var total int

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		_ = json.NewDecoder(r.Body).Decode(&url)

		err := db.QueryRow("SELECT COUNT(*) from sp_urls WHERE url = ?", url.URL).Scan(&total)
		if err != nil {
			common.Message(w, "Cant count urls!")
		}

		if total == 0 {
			stmt, err := db.Prepare("INSERT INTO sp_urls(name,url,return_code,content,check_interval) values(?,?,?,?,?)")
			if err != nil {
				common.Message(w, "Cant prepare insert url!")
			}

			_, err = stmt.Exec(url.Name, url.URL, url.ReturnCode, url.Content, url.CheckInterval)
			if err != nil {
				common.Message(w, "Cant insert url!")
			} else {
				common.Message(w, "Url "+url.URL+" added!")
			}
		} else {
			common.Message(w, "Url already in the system!")
		}
	} else {
		common.Message(w, "Invalid token!")
	}
}

func DeleteUrl(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		params := mux.Vars(r)

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		stmt, err := db.Prepare("DELETE FROM sp_urls WHERE id = ?")
		if err != nil {
			common.Message(w, "Cant prepare delete url!")
		}

		_, err = stmt.Exec(params["id"])
		if err != nil {
			common.Message(w, "Cant delete url!")
		} else {
			common.Message(w, "Url "+params["id"]+" deleted!")
		}
	} else {
		common.Message(w, "Invalid token!")
	}
}

func ReturnURLInfo(IDUrl int) Url {

	var url Url

	db, errDB := common.DBConnection()
	defer db.Close()
	if errDB != nil {
		log.Println("Cant connect to server host!")
	}

	err := db.QueryRow("SELECT id,name,url,return_code,content,check_interval FROM sp_urls WHERE id = ?", IDUrl).Scan(&url.ID, &url.Name, &url.URL, &url.ReturnCode, &url.Content, &url.CheckInterval)

	if err != nil {
		log.Printf("Cant get url info!")
	}

	return url
}

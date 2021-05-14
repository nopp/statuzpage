package incidents

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"statuzpage-api/common"
	"statuzpage-api/urls"

	"github.com/gorilla/mux"
)

type Incident struct {
	ID         int
	IDUrl      int
	StartedAt  string
	FinishedAt sql.NullString
	GroupName  string
	UrlName    string
	Message    string
}

type IncidentSuport struct {
	UrlName    string
	Url        string
	StartedAt  string
	FinishedAt sql.NullString
	Message    string
}

var incidents []Incident

func CreateIncident(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		var incident Incident
		var total int

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		_ = json.NewDecoder(r.Body).Decode(&incident)

		urlInfo := urls.ReturnURLInfo(incident.IDUrl)

		err := db.QueryRow("SELECT COUNT(*) from sp_incidents WHERE idUrl = ? AND finishedat IS NULL", incident.IDUrl).Scan(&total)
		if err != nil {
			common.Message(w, "Cant count incidents!")
		}

		if total == 0 {
			stmt, err := db.Prepare("INSERT INTO sp_incidents(idUrl,startedat,message) values(?,?,?)")
			if err != nil {
				common.Message(w, "Cant prepare insert incident!")
			}
			res, err := stmt.Exec(incident.IDUrl, incident.StartedAt, incident.Message)
			if err != nil {
				common.Message(w, urlInfo.Name+" cant insert incident!")
			} else {

				lastID, _ := res.LastInsertId()
				incident.ID = int(lastID)

				json.NewEncoder(w).Encode(incident)
				common.Message(w, urlInfo.Name+" incident created!")
			}

		} else {
			common.Message(w, urlInfo.Name+" incident already in database!")
		}
	} else {
		common.Message(w, "Invalid token!")
	}
}

func CloseIncident(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		params := mux.Vars(r)

		var incident Incident
		var total int

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		_ = json.NewDecoder(r.Body).Decode(&incident)

		err := db.QueryRow("SELECT COUNT(*) from sp_incidents WHERE id = ? AND finishedat IS NULL", params["id"]).Scan(&total)
		if err != nil {
			common.Message(w, "Cant count incidents!")
		}

		if total != 0 {
			stmt, err := db.Prepare("UPDATE sp_incidents SET finishedat = ? WHERE id = ?")
			if err != nil {
				common.Message(w, "Cant prepare update incident!")
			}
			res, err := stmt.Exec(incident.FinishedAt.String, params["id"])
			if err != nil {
				common.Message(w, "Cant update incident!")
			}
			incident.FinishedAt.Valid = true

			lastID, _ := res.LastInsertId()
			incident.ID = int(lastID)

			json.NewEncoder(w).Encode(incident)
			common.Message(w, "Incident "+params["id"]+" closed!")
		} else {
			common.Message(w, "Incident is not open!")
		}
	} else {
		common.Message(w, "Invalid token!")
	}
}

func GetIncidents(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		var incidents []IncidentSuport
		var incident IncidentSuport

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}
		rows, err := db.Query("SELECT url.name,url.url,inci.startedat,inci.message FROM sp_urls url JOIN sp_incidents inci on inci.idUrl = url.id AND inci.finishedat is NULL")
		if err != nil {
			common.Message(w, "Cant get incidents from systems!")
		}

		for rows.Next() {
			err := rows.Scan(&incident.UrlName, &incident.Url, &incident.StartedAt, &incident.Message)
			if err != nil {
				common.Message(w, "Cant return incident informations!")
			}
			incidents = append(incidents, incident)
		}

		json.NewEncoder(w).Encode(incidents)
	} else {
		common.Message(w, "Invalid token!")
	}
}

func GetIncidentsClosed(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		var incidents []IncidentSuport
		var incident IncidentSuport

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}
		rows, err := db.Query("SELECT i.startedat,i.finishedat,i.message,u.name from sp_incidents i, sp_urls u WHERE i.idUrl = u.id AND i.finishedat IS NOT NULL ORDER by i.finishedat DESC")
		if err != nil {
			common.Message(w, "Cant get incidents closed from systems!")
		}

		for rows.Next() {
			err := rows.Scan(&incident.StartedAt, &incident.FinishedAt, &incident.Message, &incident.UrlName)
			if err != nil {
				common.Message(w, "Cant return incidents closed informations!")
			}
			incidents = append(incidents, incident)
		}

		json.NewEncoder(w).Encode(incidents)
	} else {
		common.Message(w, "Invalid token!")
	}
}

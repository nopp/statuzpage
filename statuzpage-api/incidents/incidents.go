package incidents

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"statuzpage-api/common"
	"statuzpage-api/urls"
	"strconv"

	"github.com/gorilla/mux"
)

type Incident struct {
	ID         int
	IDGroup    int
	IDUrl      int
	StartedAt  string
	FinishedAt sql.NullString
	Message    string
	GroupName  string
	UrlName    string
}

var incidents []Incident

// Incident
func GetIncidents(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		var incidents []Incident
		var incident Incident
		params := mux.Vars(r)

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}
		if params["limit"] == "0" {
			rows, err := db.Query("SELECT i.id,i.idGroup,i.idUrl,i.startedat,i.finishedat,i.message,g.name,u.name FROM sp_incidents i, sp_groups g, sp_urls u WHERE i.idGroup = g.id AND i.idUrl = u.id ORDER by i.id DESC")
			if err != nil {
				common.Message(w, "Cant get incidents from systems!")
			}

			for rows.Next() {
				err := rows.Scan(&incident.ID, &incident.IDGroup, &incident.IDUrl, &incident.StartedAt, &incident.FinishedAt, &incident.Message, &incident.GroupName, &incident.UrlName)
				if err != nil {
					common.Message(w, "Cant return incident informations!")
				}
				incidents = append(incidents, incident)
			}

			json.NewEncoder(w).Encode(incidents)
		} else {
			rows, err := db.Query("SELECT i.id,i.idGroup,i.idUrl,i.startedat,i.finishedat,i.message,g.name,u.name FROM sp_incidents i, sp_groups g, sp_urls u WHERE i.idGroup = g.id AND i.idUrl = u.id ORDER by i.id DESC LIMIT ?", params["limit"])
			if err != nil {
				common.Message(w, "Cant get incidents from systems!")
			}

			for rows.Next() {
				err := rows.Scan(&incident.ID, &incident.IDGroup, &incident.IDUrl, &incident.StartedAt, &incident.FinishedAt, &incident.Message, &incident.GroupName, &incident.UrlName)
				if err != nil {
					common.Message(w, "Cant return incident informations!")
				}
				incidents = append(incidents, incident)
			}

			json.NewEncoder(w).Encode(incidents)
		}
	} else {
		common.Message(w, "Invalid token!")
	}
}

func GetIncidentsByIdGroup(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		var incidents []Incident
		var incident Incident
		params := mux.Vars(r)

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}
		if params["limit"] == "0" {
			rows, err := db.Query("SELECT i.id,i.idGroup,i.idUrl,i.startedat,i.finishedat,i.message,g.name,u.name FROM sp_incidents i, sp_groups g, sp_urls u WHERE i.idGroup = g.id AND i.idUrl = u.id AND i.idGroup = ? ORDER by i.id DESC", params["idgroup"])
			if err != nil {
				common.Message(w, "Cant get incidents from systems!")
			}
			for rows.Next() {
				err := rows.Scan(&incident.ID, &incident.IDGroup, &incident.IDUrl, &incident.StartedAt, &incident.FinishedAt, &incident.Message, &incident.GroupName, &incident.UrlName)
				if err != nil {
					common.Message(w, "Cant return incident informations!")
				}
				incidents = append(incidents, incident)
			}

			json.NewEncoder(w).Encode(incidents)
		} else {
			rows, err := db.Query("SELECT i.id,i.idGroup,i.idUrl,i.startedat,i.finishedat,i.message,g.name,u.name FROM sp_incidents i, sp_groups g, sp_urls u WHERE i.idGroup = g.id AND i.idUrl = u.id AND i.idGroup = ? ORDER by i.id DESC LIMIT ?", params["idgroup"], params["limit"])
			if err != nil {
				common.Message(w, "Cant get incidents from systems!")
			}
			for rows.Next() {
				err := rows.Scan(&incident.ID, &incident.IDGroup, &incident.IDUrl, &incident.StartedAt, &incident.FinishedAt, &incident.Message, &incident.GroupName, &incident.UrlName)
				if err != nil {
					common.Message(w, "Cant return incident informations!")
				}

				incidents = append(incidents, incident)
			}

			json.NewEncoder(w).Encode(incidents)
		}
	} else {
		common.Message(w, "Invalid token!")
	}
}

func GetIncidentsByIdGroupMonthYear(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		var incidents []Incident
		var incident Incident
		params := mux.Vars(r)

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}
		rows, err := db.Query("SELECT i.id,i.idGroup,i.idUrl,i.startedat,i.finishedat,i.message,g.name,u.name FROM sp_incidents i, sp_groups g, sp_urls u WHERE i.idGroup = g.id AND i.idUrl = u.id AND i.idGroup = ? AND MONTH(i.startedat) = ? AND YEAR(i.startedat) = ? ORDER by i.id DESC", params["idgroup"], params["month"], params["year"])
		if err != nil {
			common.Message(w, "Cant get incidents from systemsX!")
		}
		for rows.Next() {
			err := rows.Scan(&incident.ID, &incident.IDGroup, &incident.IDUrl, &incident.StartedAt, &incident.FinishedAt, &incident.Message, &incident.GroupName, &incident.UrlName)
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

func GetIncident(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		params := mux.Vars(r)
		var incident Incident

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		err := db.QueryRow("SELECT id,idGroup,idUrl,startedat,finishedat,message from sp_incidents WHERE id = ?", params["id"]).Scan(&incident.ID, &incident.IDGroup, &incident.IDUrl, &incident.StartedAt, &incident.FinishedAt, &incident.Message)
		if err != nil {
			common.Message(w, "Cant get incident from systems!")
		}

		json.NewEncoder(w).Encode(incident)
	} else {
		common.Message(w, "Invalid token!")
	}
}

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

		err := db.QueryRow("SELECT COUNT(*) from sp_incidents WHERE idGroup = ? AND idUrl = ? AND finishedat IS NULL", incident.IDGroup, incident.IDUrl).Scan(&total)
		if err != nil {
			common.Message(w, "Cant count incidents!")
		}

		if total == 0 {
			stmt, err := db.Prepare("INSERT INTO sp_incidents(idGroup,idUrl,startedat,message) values(?,?,?,?)")
			if err != nil {
				common.Message(w, "Cant prepare insert incident!")
			}
			res, err := stmt.Exec(int(incident.IDGroup), incident.IDUrl, incident.StartedAt, incident.Message)
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
			stmt, err := db.Prepare("UPDATE sp_incidents SET finishedat = ?, message = ? WHERE id = ?")
			if err != nil {
				common.Message(w, "Cant prepare update incident!")
			}
			res, err := stmt.Exec(incident.FinishedAt.String, incident.Message, params["id"])
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

func DeleteIncident(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		params := mux.Vars(r)

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		stmt, err := db.Prepare("DELETE FROM sp_incidents WHERE id = ?")
		if err != nil {
			common.Message(w, "Cant prepare delete incidents!")
		}

		_, err = stmt.Exec(params["id"])
		if err != nil {
			common.Message(w, "Cant delete incident!")
		}

		common.Message(w, "Incident "+params["id"]+" deleted!")
	} else {
		common.Message(w, "Invalid token!")
	}
}

func GetTotalIncidentsOpen(w http.ResponseWriter, r *http.Request) {

	if common.CheckToken(r.Header.Get("statuzpage-token")) {

		var total int

		db, errDB := common.DBConnection()
		defer db.Close()
		if errDB != nil {
			common.Message(w, "Cant connect to server host!")
		}

		err := db.QueryRow("SELECT count(*) FROM sp_incidents i, sp_groups g WHERE i.idGroup = g.id AND i.finishedat IS NULL").Scan(&total)
		if err != nil {
			common.Message(w, "Cant count open incidents !")
		}
		common.Message(w, strconv.Itoa(total))
	} else {
		common.Message(w, "Invalid token!")
	}
}

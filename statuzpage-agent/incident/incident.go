package incident

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"statuzpage-agent/configuration"
	"statuzpage-agent/db"
	"strconv"
)

type incident struct {
	ID         int            `json:"id"`
	IDGroup    int            `json:"idgroup"`
	IDUrl      int            `json:"idurl,omitempty"`
	StartedAt  string         `json:"startedat,omitempty"`
	FinishedAt sql.NullString `json:"finishedat,omitempty"`
	Message    string         `json:"message"`
}

// CreateIncident responsible for create a new incident
func CreateIncident(idGroup, idURL int, message, AppName, GroupName, startedAt string) {

	var incident incident

	config := configuration.LoadConfiguration()

	transCfg := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transCfg}

	incident.IDGroup = idGroup
	incident.IDUrl = idURL
	incident.StartedAt = startedAt
	incident.Message = message
	incidentJSON, _ := json.Marshal(incident)

	req, err := http.NewRequest("POST", "http://"+config.StatuzpageAPI+"/incident", bytes.NewBuffer([]byte(incidentJSON)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("statuzpage-token", config.Token)
	if err != nil {
		log.Println(AppName, "can't create incident!")
	}

	resp, errDO := client.Do(req)
	if errDO != nil {
		log.Println(errDO)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(GroupName + " (" + AppName + ")" + message + " " + string(body))
}

// CloseIncident responsible for close incident opened
func CloseIncident(id int, finishedAt, message, AppName, GroupName string) {

	var incident incident

	config := configuration.LoadConfiguration()

	transCfg := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transCfg}

	incident.ID = id
	incident.FinishedAt.String = finishedAt
	incident.FinishedAt.Valid = true
	incident.Message = message
	incidentJSON, _ := json.Marshal(incident)

	req, err := http.NewRequest("POST", "http://"+config.StatuzpageAPI+"/incident/"+strconv.Itoa(id), bytes.NewBuffer(incidentJSON))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("statuzpage-token", config.Token)
	if err != nil {
		log.Println(AppName + " cant close incident!")
	}

	resp, errDO := client.Do(req)
	if errDO != nil {
		log.Println(errDO)
	}
	defer resp.Body.Close()
	log.Println(GroupName + " (" + AppName + ") incident closed!")
}

// IsOpen verify if incident was opened
func IsOpen(IDGroup, IDUrl int) bool {

	var total int

	db, errDB := db.DBConnection()
	defer db.Close()
	if errDB != nil {
		log.Println("Cant connect to server host!")
	}

	err := db.QueryRow("SELECT COUNT(*) from sp_incidents WHERE idGroup = ? AND idUrl = ? AND finishedat IS NULL", IDGroup, IDUrl).Scan(&total)
	if err != nil {
		log.Println("Cant count incidents!")
	}

	if total == 0 {
		return false
	} else {
		return true
	}
}

// ReturnIDIncidentOpen return id from opened incidente
func ReturnIDIncidentOpen(IDGroup, IDUrl int) int {

	var id int

	db, errDB := db.DBConnection()
	defer db.Close()
	if errDB != nil {
		log.Println("Cant connect to server host!")
	}

	err := db.QueryRow("SELECT id from sp_incidents WHERE idGroup = ? AND idUrl = ? AND finishedat IS NULL", IDGroup, IDUrl).Scan(&id)
	if err != nil {
		log.Println("Cant count incidents!")
	}

	return id
}

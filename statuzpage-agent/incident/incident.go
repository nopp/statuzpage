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
	IDUrl      int            `json:"idurl,omitempty"`
	StartedAt  string         `json:"startedat,omitempty"`
	FinishedAt sql.NullString `json:"finishedat,omitempty"`
	Message    string         `json:"message"`
}

type msg struct {
	Message string `json:"message"`
}

// CreateIncident responsible for create a new incident
func CreateIncident(idURL int, message, AppName, startedAt string) {

	var incident incident
	var msgJson msg

	config := configuration.LoadConfiguration()

	transCfg := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transCfg}

	incident.IDUrl = idURL
	incident.StartedAt = startedAt
	incident.Message = message
	incidentJSON, _ := json.Marshal(incident)

	req, err := http.NewRequest("POST", "http://"+config.StatuzpageAPI+"/incident", bytes.NewBuffer([]byte(incidentJSON)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("statuzpage-token", config.Token)
	if err != nil {
		log.Println(AppName, " can't create incident!")
	}

	resp, errDO := client.Do(req)
	if errDO != nil {
		log.Println(errDO)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &msgJson)
	log.Println(AppName + " " + message + " " + msgJson.Message)
}

// CloseIncident responsible for close incident opened
func CloseIncident(id int, finishedAt, AppName string) {

	var incident incident

	config := configuration.LoadConfiguration()

	transCfg := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: transCfg}

	incident.ID = id
	incident.FinishedAt.String = finishedAt
	incident.FinishedAt.Valid = true
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
	log.Println(AppName + " incident closed!")
}

// IsOpen verify if incident was opened
func IsOpen(IDUrl int) bool {

	var total int

	db, errDB := db.DBConnection()
	defer db.Close()
	if errDB != nil {
		log.Println("Cant connect to server host!")
	}

	err := db.QueryRow("SELECT COUNT(*) from sp_incidents WHERE idUrl = ? AND finishedat IS NULL", IDUrl).Scan(&total)
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
func ReturnIDIncidentOpen(IDUrl int) int {

	var id int

	db, errDB := db.DBConnection()
	defer db.Close()
	if errDB != nil {
		log.Println("Cant connect to server host!")
	}

	err := db.QueryRow("SELECT id from sp_incidents WHERE idUrl = ? AND finishedat IS NULL", IDUrl).Scan(&id)
	if err != nil {
		log.Println("Cant count incidents!")
	}

	return id
}

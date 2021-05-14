package agent

// Package agent
// I created this package to solve "Import Cycles"
// I don't know if this is the best way, but, resolved.

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"statuzpage-agent/db"
	"statuzpage-agent/endpoint"
	"statuzpage-agent/incident"
	"time"

	"github.com/jasonlvhit/gocron"
)

// ReturnAllUrls Check all urls
func ReturnAllUrls() {

	var urlStruct endpoint.Info

	db, errDB := db.DBConnection()
	defer db.Close()
	if errDB != nil {
		log.Println("Cant connect to server host!")
	}

	urls, err := endpoint.ReturnUrls()
	if err != nil {
		log.Println("Cant connect to server host!")
	}
	for urls.Next() {
		err := urls.Scan(&urlStruct.ID, &urlStruct.Name, &urlStruct.URL, &urlStruct.ReturnCode, &urlStruct.Content, &urlStruct.CheckInterval)
		if err != nil {
			log.Print(err)
		} else {
			gocron.Every(urlStruct.CheckInterval).Seconds().Do(Check, urlStruct.ID, urlStruct.Name, urlStruct.URL)
		}
	}
	<-gocron.Start()
}

// Check responsible for health of url(endpoint)
func Check(IDUrl int, AppName, url string) {

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	result, err := client.Get(url)
	if err != nil {
		incident.CreateIncident(IDUrl, "unreachable!", AppName, time.Now().Format("2006-01-02 15:04:05"))
	} else {
		if result.StatusCode == 200 {
			var url = endpoint.ReturnURLInfo(IDUrl)
			// If url have content to check
			if url.Content.Valid {
				content, errContent := ioutil.ReadAll(result.Body)
				if errContent != nil {
					log.Println("Read content problem! " + errContent.Error())
				} else {
					matched, errMatch := regexp.MatchString(url.Content.String, string(content))
					if errMatch != nil {
						log.Println("Content match problem! " + errMatch.Error())
					} else {
						if matched {
							if incident.IsOpen(IDUrl) {
								incident.CloseIncident(incident.ReturnIDIncidentOpen(IDUrl), time.Now().Format("2006-01-02 15:04:05"), AppName)
							} else {
								log.Println(AppName + " operational and content matched!")
							}
						} else {
							incident.CreateIncident(IDUrl, "operational, but content doesn't match!", AppName, time.Now().Format("2006-01-02 15:04:05"))
						}
					}
				}
			} else {
				if incident.IsOpen(IDUrl) {
					incident.CloseIncident(incident.ReturnIDIncidentOpen(IDUrl), time.Now().Format("2006-01-02 15:04:05"), AppName)
				} else {
					log.Println(AppName + " operational!")
				}
			}
		} else {
			incident.CreateIncident(IDUrl, "with problem!", AppName, time.Now().Format("2006-01-02 15:04:05"))
		}
		defer result.Body.Close()
	}
}

package eg

// Package endpointGroups (eg)
// I created this package to solve "Import Cycles"
// I don't know if this is the best way, but, resolved.

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"statuzpage-agent/db"
	"statuzpage-agent/endpoint"
	"statuzpage-agent/group"
	"statuzpage-agent/incident"
	"time"

	"github.com/jasonlvhit/gocron"
)

// ReturnUrlsAllGroups Check all urls from all groups
func ReturnUrlsAllGroups() {

	var groupInfo group.Group
	var urlStruct endpoint.Info

	db, errDB := db.DBConnection()
	defer db.Close()
	if errDB != nil {
		log.Println("Cant connect to server host!")
	}

	result, err := db.Query("SELECT id,name FROM sp_groups")
	if err != nil {
		log.Println(err)
	} else {
		for result.Next() {
			err := result.Scan(&groupInfo.ID, &groupInfo.Name)
			if err != nil {
				fmt.Println(err)
			}
			urls, err := endpoint.ReturnUrlsByIDGroup(groupInfo.ID)
			for urls.Next() {
				err := urls.Scan(&urlStruct.ID, &urlStruct.Name, &urlStruct.URL, &urlStruct.ReturnCode, &urlStruct.Content, &urlStruct.CheckInterval)
				if err != nil {
					log.Print(err)
				} else {
					gocron.Every(urlStruct.CheckInterval).Seconds().Do(Check, groupInfo.ID, urlStruct.ID, urlStruct.Name, urlStruct.URL)
				}
			}
		}
		<-gocron.Start()
	}
}

// CheckByIDGroup Check url by ID group
func CheckByIDGroup(IDGroup int) {

	var url endpoint.Info
	urls, err := endpoint.ReturnUrlsByIDGroup(IDGroup)
	if err != nil {
		fmt.Println(err)
		log.Println("Cant return app servers!")
	} else {
		for urls.Next() {
			err := urls.Scan(&url.ID, &url.Name, &url.URL, &url.ReturnCode, &url.Content, &url.CheckInterval)
			fmt.Println(err)
			if err != nil {
				log.Print(err)
			} else {
				gocron.Every(url.CheckInterval).Seconds().Do(Check, IDGroup, url.ID, url.Name, url.URL)
			}
		}
		<-gocron.Start()
	}
}

// Check responsible for health of url(endpoint)
func Check(IDGroup, IDUrl int, AppName, url string) {

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	groupName := group.ReturnGroupName(IDGroup)

	if groupName == "" {
		log.Println("Cant get group name!")
	} else {
		result, err := client.Get(url)
		if err != nil {
			incident.CreateIncident(IDGroup, IDUrl, "unreachable!", AppName, groupName, time.Now().Format("2006-01-02 15:04:05"))
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
								if incident.IsOpen(IDGroup, IDUrl) {
									incident.CloseIncident(incident.ReturnIDIncidentOpen(IDGroup, IDUrl), time.Now().Format("2006-01-02 15:04:05"), "solved!", AppName, groupName)
								} else {
									log.Println(groupName + " (" + AppName + ") operational and content matched!")
								}
							} else {
								incident.CreateIncident(IDGroup, IDUrl, "operational, but content doesn't match!", AppName, groupName, time.Now().Format("2006-01-02 15:04:05"))
							}
						}
					}
				} else {
					if incident.IsOpen(IDGroup, IDUrl) {
						incident.CloseIncident(incident.ReturnIDIncidentOpen(IDGroup, IDUrl), time.Now().Format("2006-01-02 15:04:05"), "solved!", AppName, groupName)
					} else {
						log.Println(groupName + " (" + AppName + ") operational!")
					}
				}
			} else {
				incident.CreateIncident(IDGroup, IDUrl, "with problem!", AppName, groupName, time.Now().Format("2006-01-02 15:04:05"))
			}
			defer result.Body.Close()
		}
	}

}

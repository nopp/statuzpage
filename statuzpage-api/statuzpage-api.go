package main

import (
	"log"
	"net/http"
	"statuzpage-api/common"
	"statuzpage-api/groups"
	"statuzpage-api/incidents"
	"statuzpage-api/urls"
	"statuzpage-api/users"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	// User methods
	router.HandleFunc("/users", users.GetUsers).Methods("GET")
	router.HandleFunc("/user/{id}", users.GetUser).Methods("GET")
	router.HandleFunc("/user", users.CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", users.DeleteUser).Methods("DELETE")
	// URL methods
	router.HandleFunc("/urls", urls.GetUrls).Methods("GET")
	router.HandleFunc("/url/{id}", urls.GetUrl).Methods("GET")
	router.HandleFunc("/url", urls.CreateUrl).Methods("POST")
	router.HandleFunc("/url/{id}", urls.DeleteUrl).Methods("DELETE")
	// Group methods
	router.HandleFunc("/groups", groups.GetGroups).Methods("GET")
	router.HandleFunc("/group/{id}", groups.GetGroup).Methods("GET")
	router.HandleFunc("/group", groups.CreateGroup).Methods("POST")
	router.HandleFunc("/group/{id}", groups.DeleteGroup).Methods("DELETE")
	// Incident methods
	router.HandleFunc("/incidents/limit/{limit}", incidents.GetIncidents).Methods("GET")
	router.HandleFunc("/incidents/{idgroup}/limit/{limit}", incidents.GetIncidentsByIdGroup).Methods("GET")
	router.HandleFunc("/incidents/{idgroup}/{month}/{year}", incidents.GetIncidentsByIdGroupMonthYear).Methods("GET")
	router.HandleFunc("/incident/{id}", incidents.GetIncident).Methods("GET")
	router.HandleFunc("/incident", incidents.CreateIncident).Methods("POST")
	router.HandleFunc("/incident/{id}", incidents.CloseIncident).Methods("POST")
	router.HandleFunc("/incident/{id}", incidents.DeleteIncident).Methods("DELETE")
	// Tools
	router.HandleFunc("/totalOpen", incidents.GetTotalIncidentsOpen).Methods("GET")

	config := common.LoadConfiguration()
	log.Fatal(http.ListenAndServe(config.HostPort, router))
}

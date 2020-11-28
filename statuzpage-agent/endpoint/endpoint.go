package endpoint

import (
	"database/sql"
	"log"
	"statuzpage-agent/db"
)

// Info public struc
type Info struct {
	ID            int
	IDGroup       int
	Name          string
	URL           string
	ReturnCode    string
	Content       sql.NullString
	CheckInterval uint64
}

// ReturnURLInfo Return information about url
func ReturnURLInfo(IDUrl int) Info {

	var url Info

	db, errDB := db.DBConnection()
	defer db.Close()
	if errDB != nil {
		log.Println("Cant connect to server host!")
	}

	err := db.QueryRow("SELECT id,idGroup,name,url,return_code,content,check_interval FROM sp_urls WHERE id = ?", IDUrl).Scan(&url.ID, &url.IDGroup, &url.Name, &url.URL, &url.ReturnCode, &url.Content, &url.CheckInterval)

	if err != nil {
		log.Printf("Cant get url info!")
	}

	return url
}

// ReturnUrlsByIDGroup Return url from ID Group
func ReturnUrlsByIDGroup(IDGroup int) (*sql.Rows, error) {

	db, errDB := db.DBConnection()
	defer db.Close()
	if errDB != nil {
		log.Println("Cant connect to server host!")
	}

	result, err := db.Query("SELECT id,name,url,return_code,content,check_interval FROM sp_urls WHERE idGroup = ?", IDGroup)

	return result, err
}

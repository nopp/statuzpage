package main

import (
	"statuzpage-agent/agent"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	agent.ReturnAllUrls()
}

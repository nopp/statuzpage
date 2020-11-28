package main

import (
	"fmt"
	"os"
	"statuzpage-agent/eg"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func shortDur(d time.Duration) string {

	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}

	return s
}

func main() {

	options := `Usage:
	statuzpage-agent all
	statuzpage-agent group idGroup
	`
	argsWithProg := os.Args
	totalArguments := len(argsWithProg)
	if totalArguments > 1 {
		switch argsWithProg[1] {
		case "group":
			if totalArguments == 3 {
				idInt, _ := strconv.Atoi(argsWithProg[2])
				eg.CheckByIDGroup(idInt)
			} else {
				fmt.Println(options)
			}
		case "all":
			if totalArguments == 2 {
				eg.ReturnUrlsAllGroups()
			} else {
				fmt.Println(options)
			}
		default:
			fmt.Println(options)
		}
	} else {
		fmt.Println(options)
	}
}

# Statuzpage API

Requirements:
=============
	go get github.com/go-sql-driver/mysql
	go get github.com/gorilla/mux

Configurations:
===============
Change informations of config.json and copy to /etc/statuzpage-api/config.json

Obs:. Token value is used on statuzpage-token header in all actions!

Run with container:
===================
	docker build -t statuzpage-api:latest .
	docker run -d -p 8000:8000 statuzpage-api:latest
	OR
	docker run -d -p 8000:8000 -v /localDir/config.json:/etc/statuzpage-api/config.json nopp/statuzpage-api:latest
	
Run without container:
======================
	go build -o statuzpage-api .
	./statuzpage-api

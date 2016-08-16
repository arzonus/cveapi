.PHONY : build clean deps start stop restart logs

db:
	bash ./prepareDB.sh
deps:
	go get 

run:
	echo "run code"
	go install
	/opt/bin/cveapi

show:
	go install
	/opt/bin/cveapi

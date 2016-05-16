.PHONY : build clean deps start stop restart logs

restart: stop start logs

build:
	echo "build arzonus cveapi image"
	docker build -t cveapi .

deps:
	docker build -t cveapi:deps -f ./Dockerfile.deps .

clean:
	echo "remove arzonus cveapi image"
	docker rmi -f arzonus/cveapi


stop:
	echo "stop docker cveapi image"
	docker rm -vf cveapi

run:
	echo "run code"
	docker build -t cveapi .
	docker run --rm -p 3000:3000 --name cveapi cveapi /opt/bin/cveapi

test:
	echo "run code"
	docker build -t cveapi .
	docker run --rm -p 3000:3000 --name cveapi cveapi go test ./...

db:
	bash prepareDB.sh

logs:
	echo "fetch docker cveapi image"
	docker logs -f cveapi
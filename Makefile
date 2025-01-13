build-project:
	cd project && go get . && go build .

run-project:
	cd project && ./project

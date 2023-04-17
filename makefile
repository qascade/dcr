build: 
	go build -o bin/dcr
tidy: 
	go mod tidy
vendor: 
	go mod vendor
fmt: 
	gofmt -s -w . 

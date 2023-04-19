build: 
	go build -o bin/dcr
tidy: 
	go mod tidy
vendor: 
	go mod vendor
fmt: 
	go fmt ./...
test: build
	go test ./...
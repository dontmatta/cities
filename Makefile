build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/cities cities/get.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/states states/get.go
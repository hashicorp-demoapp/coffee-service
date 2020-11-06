build_linux:
	CGO_ENABLED=0 GOOS=linux go build -o ./bin/coffee-service
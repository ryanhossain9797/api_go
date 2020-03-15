dev:
	go run main.go
build_arm:
	env GOOS=linux GOARCH=arm go build -o api_arm
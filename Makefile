run:
	go run -race cmd/rester/main.go

build:
	go build -o rester cmd/rester/main.go
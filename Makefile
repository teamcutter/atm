run:
	go run ./main.go

serve:
	go run ./cmd/server/main.go

client:
	go run ./cmd/client/main.go

build:
	go build -o ./bin/atm && ./bin/atm

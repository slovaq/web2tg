build:
	go build -o bin/main cmd/*.go
run:
	go run cmd/*.go
clean:
	rm bin/*

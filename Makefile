build:
	go build -o serve cmd/app.go
run: 
	./serve	
all: build run	

all: k12-booter

k12-booter: *.go
	go build

clean:
	rm -f k12-booter

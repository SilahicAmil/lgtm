build:
	go build -o blame-me-later ./cmd

run:
	go run ./cmd

clean:
	rm -f blame-me-later
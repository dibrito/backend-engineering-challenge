build: clean
	go build -o calculator-cli .

clean:
	rm -rf calculator-cli
	rm -rf result.txt

run: build
	./calculator-cli --input_file ./events.json --window 10

test:
	go test -v -cover -short ./...
	make clean

PHONY: build clean run test
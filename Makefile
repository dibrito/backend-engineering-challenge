build: clean
	go build -o calculator-cli .

clean:
	rm -rf calculator-cli
	rm -rf result.txt

run: build
	./calculator-cli --input_file events.json --window 10

PHONY: build clean run
build: clean
	go build -o calculator-cli .

clean:
	rm -rf calculator-cli
	rm -rf result.txt

run: build
	./calculator-cli --input_file ./events.json --window 10

test:
	go test -v -cover -short ./... -count=1
	make clean

benchclean:
	rm -rf sma.bench
	rm -rf fifo.bench

benchsma:
	go test ./cmd -run=^$$ -benchmem -bench=^BenchmarkSMA$$ -count=10 > sma.bench

benchfifo:
	go test ./cmd -run=^$$ -benchmem -bench=^BenchmarkFIFOSMA$$ -count=10 > fifo.bench

PHONY: build clean run test benchsma benchfifo benchclean
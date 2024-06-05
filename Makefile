build: clean
	go build -o calculator-cli .

clean:
	rm -rf calculator-cli
	rm -rf result.txt
	rm -rf *.bench
	rm -rf *.pprof

run: build
	./calculator-cli --input_file ./events.json --window_size 10

test: clean
	go test -v -cover -short ./... -count=1

benchsma:
	go test ./cmd -run=^$$ -benchmem -bench=^BenchmarkSMA$$ -count=10 > sma.bench

benchcpufifo:
	go test ./cmd -run=^$$ -benchmem -bench=^BenchmarkFIFOSMA$$ -cpuprofile=cpufifo.pprof -count=10 > cpufifo.bench

benchmemfifo:
	go test ./cmd -run=^$$ -benchmem -bench=^BenchmarkFIFOSMA$$ -memprofile=memfif.pprof -count=10 > memfifo.bench

benchcpubufffifo:
	go test ./cmd -run=^$$ -benchmem -bench=^BenchmarkBuffFIFOSMA$$ -cpuprofile=cpubufffifo.pprof -count=10 > cpubufffifo.bench

benchmembufffifo:
	go test ./cmd -run=^$$ -benchmem -bench=^BenchmarkBuffFIFOSMA$$ -memprofile=membufffifo.pprof -count=10 > membufffifo.bench

PHONY: build clean run test benchsma benchfifo benchclean benchfifomepprof
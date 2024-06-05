# Calculator

Calculator is a CLI writen in Go for computing math operations over translation metrics.

# Overview

Calculator provides a simple interface to apply math operations over a specifc data scenario.

Calculator provides:

* Simple Moving Average(sma)
* more functionalities TBD

# Data scenario

Considering a translation event in the format:

```json
{
	"timestamp": "2018-12-26 18:12:19.903159",
	"translation_id": "5aa5b2f39f7254a75aa4",
	"source_language": "en",
	"target_language": "fr",
	"client_name": "airliberty",
	"event_name": "translation_delivered",
	"duration": 20,
	"nr_words": 100
}
```

When interested in calculating, for every minute, a moving average(sma) of the translations delivery time for the last X minutes, you can call calculator as bellow:

```bash
calculator --input_file events.json --window_size 10
```

The output will be write in a resukt.txt file, and should look like:

````txt
{"date":"2018-12-26 18:11:00","average_delivery_time":0}
{"date":"2018-12-26 18:12:00","average_delivery_time":20}
{"date":"2018-12-26 18:13:00","average_delivery_time":20}
{"date":"2018-12-26 18:14:00","average_delivery_time":20}
{"date":"2018-12-26 18:15:00","average_delivery_time":20}
{"date":"2018-12-26 18:16:00","average_delivery_time":25.5}
{"date":"2018-12-26 18:17:00","average_delivery_time":25.5}
{"date":"2018-12-26 18:18:00","average_delivery_time":25.5}
{"date":"2018-12-26 18:19:00","average_delivery_time":25.5}
{"date":"2018-12-26 18:20:00","average_delivery_time":25.5}
{"date":"2018-12-26 18:21:00","average_delivery_time":25.5}
{"date":"2018-12-26 18:22:00","average_delivery_time":31}
{"date":"2018-12-26 18:23:00","average_delivery_time":31}
{"date":"2018-12-26 18:24:00","average_delivery_time":42.5}
````

# Installing

Using Calculator is easy.

Clone the repo and run:

````bash
make build
````

calculate with:

```bash
calculator --input_file events.json --window_size 10
```

# Future improvements(TODOs)

TBD

# Benchmark

The whole process of benchmark is described here:

[here](./benchmarksection.md)

Final benchmark:

[here](./finalbenchresults.md)
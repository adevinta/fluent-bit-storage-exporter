# Fluent Bit Storage Exporter

## What is this?

A Prometheus Fluent Bit exporter that serves in a Prometheus format the input and storage layer metrics already provided by Fluent Bit in a JSON format.

## How does it work?

The exporter targets the endpoint where the Fluent Bit provides the input and storage metrics, then serves the metrics parsed in a Prometheus format at a different endpoint as a web server.

By default, the exporter is set to target the URI `127.0.0.1:2020/api/v1/storage` to retrieve the metrics provided by Fluent Bit and then serve the parsed metrics in Prometheus format on port `8080`. 

## How to run the exporter?

There are different options to run the exporter, a couple of them are:

### Go run

First install dependencies of the project:

` go get ./...`

Then to make the exporter start, you can run:

`go run cmd/main.go`

If the above command is executed with no arguments it will use the default values. Also, you can provide three arguments: 
- The first one is the Fluent Bit host where the exporter will try to get the metrics (`127.0.0.1`by default)
- The second one is the Fluent Bit port (`2020` by default) 
- The third argument is the exporter port where it will serve the metrics parsed in Prometheus format.

### With a container

This project provides a Dockerfile so the exporter can be run inside a container.

To build the image, run:

`docker build . -t fluent-bit-storage-exporter`

## Metrics served by the exporter:

### Input metrics:

- fluentbit_storage_input_overlimit
- fluentbit_storage_input_mem_bytes
- fluentbit_storage_input_limit_bytes
- fluentbit_storage_input_chunks
- fluentbit_storage_input_chunks_fs_down
- fluentbit_storage_input_chunks_busy
- fluentbit_storage_input_busy_bytes

### Storage metrics:

- fluentbit_storage_chunks
- fluentbit_storage_chunks_mem
- fluentbit_storage_chunks_fs
- fluentbit_storage_chunks_fs_up
- fluentbit_storage_chunks_fs_down
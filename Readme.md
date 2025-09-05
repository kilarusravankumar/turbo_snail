# 🚧 Work in Progress 🚧

# TurboSnail

## About The Project

TurboSnail is a lightweight, high-performance message broker written in Go. Unlike traditional message systems like Kafka which use a FIFO (First-In, First-Out) queue, TurboSnail implements a **priority queue** for its Tracks. This allows messages to be consumed based on their assigned priority, ensuring that high-importance messages are processed first, regardless of their arrival time.

This project is designed to be a simple, embeddable, and efficient solution for scenarios where message priority is a critical requirement.

## Core Features

* **Priority-Based Messaging:** Messages with a lower priority number are consumed before messages with a higher priority number.

* **FIFO for Equal Priority:** Messages with the same priority level are guaranteed to be delivered in the order they were received.

* **Thread-Safe:** Designed for high concurrency using Go's native concurrency primitives (`sync.RWMutex`, `atomic`) to allow many producers and consumers to operate simultaneously without issues.

* **Write-Ahead Log (WAL):** TurboSnail now features a WAL to ensure message durability. Messages are stored in `gob` format. In the event of a crash, TurboSnail will look for the WAL files on disk to restore the Priority Queue messages.

* **Dynamic Track Creation:** Tracks( same concept as topics) are created on-the-fly when a message is first produced to them.

## Current Status & Limitations

This project is currently in the early stages of development. The core in-memory broker logic is functional and tested.

### TODO / Future Work


* **Consumer Groups:** Add support for consumer groups to allow multiple consumers to work together to process messages from a single Track.

* **Message Acknowledgement:** Implement an ACK/NACK mechanism for more robust message delivery guarantees.

* **Network Protocol:** Expose the broker over a network (e.g., TCP or gRPC) so it can run as a standalone service.

## Environment Variables
Create .env file in the project root.
```bash
WAL_DIR=/your/directory/for/WAL/Logs
# START_LINE_PORT is TCP port, TurboSnail listens to the START_LINE_PORT
START_LINE_PORT=7777
# FINISH_LINE_PORT is HTTP port, Msg recieving application can long poll to TurboSnail Broker. 
FINISH_LINE_PORT=7000
```


## How To Run

For Now You can run the `main.go` file to see the priority queue in action:

```bash
go run main.go 
```

Or you can run using Docker compose after creating .env, you can run below command.

```bash
docker compose up -d
```

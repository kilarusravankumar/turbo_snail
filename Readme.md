# 🚧 Work in Progress 🚧

# TurboSnail

## About The Project

TurboSnail is a lightweight, high-performance message broker written in Go. Unlike traditional message systems like Kafka which use a FIFO (First-In, First-Out) queue, TurboSnail implements a **priority queue** for its Tracks. This allows messages to be consumed based on their assigned priority, ensuring that high-importance messages are processed first, regardless of their arrival time.

This project is designed to be a simple, embeddable, and efficient solution for scenarios where message priority is a critical requirement.

## Core Features

* **Priority-Based Messaging:** Messages with a lower priority number are consumed before messages with a higher priority number.

* **FIFO for Equal Priority:** Messages with the same priority level are guaranteed to be delivered in the order they were received.

* **Thread-Safe:** Designed for high concurrency using Go's native concurrency primitives (`sync.RWMutex`, `atomic`) to allow many producers and consumers to operate simultaneously without issues.

* **Dynamic Track Creation:** Tracks( same concept as topics) are created on-the-fly when a message is first produced to them.

## Current Status & Limitations

This project is currently in the early stages of development. The core in-memory broker logic is functional and tested.

### TODO / Future Work

* **Implement Write-Ahead Log (WAL):** The most critical missing feature is persistence. Currently, the broker operates **in-memory only**. If the application crashes, all messages in the topics will be lost. A WAL will be added to ensure message durability and allow for state recovery on restart.

* **Consumer Groups:** Add support for consumer groups to allow multiple consumers to work together to process messages from a single Track.

* **Message Acknowledgement:** Implement an ACK/NACK mechanism for more robust message delivery guarantees.

* **Network Protocol:** Expose the broker over a network (e.g., TCP or gRPC) so it can run as a standalone service.

## How To Run

For Now You can run the `main.go` file to see the priority queue in action:

```bash
go run priority_broker_go

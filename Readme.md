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

* **Message Acknowledgement:** Provides an ACK/NACK mechanism for more robust message delivery guarantees.

* **Network Protocol:** Exposes the broker over both TCP (for producing) and HTTP (for consuming and acknowledging).

## Current Status & Limitations

This project is currently in the early stages of development. The core in-memory broker logic is functional and tested.

### TODO / Future Work

* **Consumer Groups:** Add support for consumer groups to allow multiple consumers to work together to process messages from a single Track.

## API

TurboSnail exposes both TCP and HTTP interfaces for message interaction.

### TCP Interface

- **Port:** The TCP server runs on the `START_LINE_PORT` (default: `7777`).
- **Functionality:** This interface is used for producing messages to a Track. To send a message, you can establish a TCP connection and send the message payload.

### HTTP Interface

- **Port:** The HTTP server runs on the `FINISH_LINE_PORT` (default: `7000`).
- **Functionality:** This interface is used for consuming messages from a Track and managing their lifecycle.

#### Endpoints

- **`GET /{track}/message`**
  - **Description:** Retrieves a message from the specified Track. This endpoint uses long polling and will wait until a message is available.
  - **Example:**
    ```bash
    curl http://localhost:7000/my-track/message
    ```

- **`POST /{track}/message/{msgId}/ack`**
  - **Description:** Acknowledges a message, confirming that it has been successfully processed.
  - **Example:**
    ```bash
    curl -X POST http://localhost:7000/my-track/message/your-message-id/ack
    ```

- **`POST /{track}/message/{msgId}/nack`**
  - **Description:** Negatively acknowledges a message, indicating that it could not be processed. The message will be requeued and delivered again.
  - **Example:**
    ```bash
    curl -X POST http://localhost:7000/my-track/message/your-message-id/nack
    ```

## Environment Variables
Create .env file in the project root.
```bash
WAL_DIR=/your/directory/for/WAL/Logs
# START_LINE_PORT is the TCP port TurboSnail listens to for incoming messages.
START_LINE_PORT=7777
# FINISH_LINE_PORT is the HTTP port for consuming messages and sending acknowledgements.
FINISH_LINE_PORT=7000
```


## How To Run

You can run the `main.go` file to start the TurboSnail broker:

```bash
go run main.go 
```

Alternatively, you can use Docker Compose after creating a `.env` file:

```bash
docker compose up -d
```

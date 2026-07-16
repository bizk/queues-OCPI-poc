# queues-OCPI-poc
PoC: queues (NATS) implementation with Mongo DB and the OCPI protocol for EV


## Development with Docker Compose

Start the API, consumer, NATS JetStream, and MongoDB:

```sh
docker compose up --build
```

The API is available at `http://localhost:8080`. Changes to Go files are picked up automatically by Air in both the API and consumer containers.

From the host, the infrastructure endpoints are:

- NATS: `nats://localhost:4222`
- NATS monitoring: `http://localhost:8222`
- MongoDB: `mongodb://localhost:27017/queues`

Inside the Compose network, the Go services should use `nats://nats:4222` and `mongodb://mongo:27017/queues`. Stop the stack with `docker compose down`; add `-v` if the MongoDB or NATS development data should also be removed.

## Dead letter queue

A Dead Letter Queue DLQ is used to store messages that for some reason cannot be processed, there might be multiple reasons, malformed message, exceeded timeout or size, the queue is full, etc.

The receiver or sender sends it to another queue to process it latter or at least investigation.


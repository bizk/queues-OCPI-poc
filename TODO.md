Producer

- [x] Init 
- [x] Create docker compose with nats
- [x] Create Rest endpoint to receive events. It should have 3 endpoints
  - [x] / as a health check
  - [x] POST /event to create a new event
  - [x] GET /event to get the events
- [x] Connect to NATS 
- [x] Produce and publish Nats message
- [x] Create job to produce random messages at random times with poison
- [ ] Create worker to do the job
- [x] Mongo
  - [x] Connect to Mongo
  - [x] Write to Mongo
- [ ] OCPI logic
- [ ] Add dead letter logic
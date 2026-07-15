Producer

- [ ] Init 
- [ ] Create docker compose with nats
- [ ] Create Rest endpoint to receive events. It should have 3 endpoints
  - [ ] / as a health check
  - [ ] POST /event to create a new event
  - [ ] GET /event to get the events
- [ ] Connect to NATS 
- [ ] Produce and publish Nats message
- [ ] Create job to produce random messages at random times with poison
- [ ] Create worker to do the job
- [ ] Mongo
  - [ ] Connect to Mongo
  - [ ] Write to Mongo
- [ ] OCPI logic
- [ ] Add dead letter logic
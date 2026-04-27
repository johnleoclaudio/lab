# Transactional Outbox Pattern
- Helps solves dual-write problem in distributed systems
- Use transactional write to both target table and the outbox table 
- A background consumer like CDC or DynamoDB streams capture these events and publish to a message queue

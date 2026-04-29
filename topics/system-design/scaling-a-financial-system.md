# Scalable Ledger

### Problem
Support high TPS backend system than is strongly consistent but highly available

### Options

##### Async Processing
1. User request for debit 
2. System queues the request and return ack response
3. Background worker consumes from the queue and process the request safely

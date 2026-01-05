# Request ID Best Practice
- always use `X-Request-ID` header
- Header always, body sometimes (return requestID in the body when encountered an error)
- use UUIDv7 as they are sortable

# System Design Checklist

### Fundamentals
- [ ] What are the read patterns?
- [ ] What are the write patterns?
- [ ] Who owns the source of truth?
- [ ] Is consistency or availability more critical?
- [ ] Singe writer or multiple writers?

### Architecture
- [ ] Sync or async?
- [ ] Do i need a queue or is a cron job enough?
- [ ] Can I separate the compute from storage?
- [ ] Stateless or stateful service?
- [ ] Contracts versioned?

### Reliability
- [ ] What happens when this fails?
- [ ] Where's the retry logic, and is it idempotent?
- [ ] Are we alerting to symptoms or root causes?
- [ ] Timeouts configured?

### Scaling
- [ ] How do reads scale?
- [ ] How do writes scale?
- [ ] Will this design hold up at 10x traffic?
- [ ] What's the hot path, and how do we optimize it?

### Observability
- [ ] Do we log what we need to debug in production?
- [ ] Can we trace a request across services?
- [ ] What metrics define "healthy"?
- [ ] Debuggable without redeploy?

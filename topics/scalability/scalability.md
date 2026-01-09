# Wed Nov 26 07:54:39 PM PST 2025
- Just got an idea how can i demonstrate ways of scaling a simple application where a participant can signup as proof of attendance.
- Rough idea would be 3 phases: 1000 users, 100k users, and 1M users.
- Technology of choice would be: Go, PostgreSQL, Docker
- For the phase 1: Monolith, go server -> postgresql
- I will have to find a way to load test this, make sure to set constraints on the resources for both db and server

```yaml
db:
  image: postgres:16-alpine
  volumes:
    - postgres_data:/var/lib/postgresql/data
volumes:
  postgres_data:
```

- `postgres_data` means a named volume. This lives in the host `/var/lib/docker/volumes/project-name_volume-name`
- the volume name `project-name_volume-name` is composed by Docker itself. If not specified, Docker will use the project directory name as prefix and second half is the volume name from the docker-compose.yml

- Done setting up server + postgresq. Server is able to write to postgresql
- To do:
  - [x] Limit server and db size - CPU, Memory
  - [x] Add monitoring - Prometheus and grafana
  - [x] Perf test 

# Sun Nov 30 11:56:25 AM PST 2025
- cadvisor to hook into docker socket to track containers' resources - CPU, Memory, Network I/O, Disk I/O.
- For application metrics, we need to instrument our application to monitor our Golden Signals - Latency, Errors, Traffic, Saturation

### Latency (How long requests take)
**What to track:**
- Request duration histogram
- P50, P95, P99 latencies

**Implementation:** wrap http handler with middleware that times each request

### Traffic (How many requests)
**What to track:**
- Total request counter
- Requests per second
- Requests by endpoint

**Implementation:** counter that increments on each request

### Errors (How many requests fail)
**What to track:**
- Error rate by type (4xx vs 5xx)
- Error counter by endpoint
- Success rate percentage

**Implementation:** counter that increments based on response status code

### Saturation (How "full" your service is)
**What to track:**
- Database connection pool usage
- Request queue depth
- Memory/CPU approaching limits

**Implementation:** Gauge metrics for current utilization

---

- Instrument Go server -> expose endpoint for Prometheus to scrape

### Prometheus Notes:
- You define metric ("the thing"), labels/dimension ("characteristics of the thing")
- e.g. 
- Prometheus Best Practices
     ✅ Do This
     1. Keep label cardinality low (<100 unique combinations per metric)
     2. Use standard metric types (Counter, Histogram, Gauge, Summary)
     3. Use promauto package (automatic registration)
     4. Pre-define label values if possible
     5. Monitor your metrics with /metrics endpoint size

     ❌ Don't Do This
     1. Don't use user IDs, emails, IPs as labels
     2. Don't use timestamps as labels
     3. Don't create metrics dynamically per request
     4. Don't use unbounded label values



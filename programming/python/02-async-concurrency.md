# Async/Await and Concurrency for Backend

Time: 1.5 hours | Focus: Why async matters and how to use it correctly

## Why Async for Backend?

**The Problem**: Traditional synchronous code blocks on I/O operations (database, HTTP requests, file access). Your CPU sits idle waiting.

```python
# Synchronous - Each operation blocks
def get_user_data(user_id: int):
    user = db.query(user_id)          # Wait 50ms
    posts = api.get_posts(user_id)    # Wait 100ms
    profile = api.get_profile(user_id) # Wait 75ms
    return combine(user, posts, profile)
    # Total: 225ms
```

```python
# Asynchronous - Operations run concurrently
async def get_user_data(user_id: int):
    user, posts, profile = await asyncio.gather(
        db.query(user_id),              # All three run
        api.get_posts(user_id),         # at the same
        api.get_profile(user_id)        # time
    )
    return combine(user, posts, profile)
    # Total: 100ms (longest operation)
```

**When to use async**:
- High I/O workloads (APIs, databases, file operations)
- Need to handle many concurrent connections (10k+ users)
- Modern Python web frameworks (FastAPI requires it)

**When NOT to use async**:
- CPU-bound tasks (use multiprocessing instead)
- Simple scripts that run once
- Working with libraries that don't support async

## Async Basics

### Defining Async Functions

```python
import asyncio

# Regular function
def fetch_user(user_id: int) -> dict:
    return {"id": user_id}

# Async function (coroutine)
async def fetch_user_async(user_id: int) -> dict:
    await asyncio.sleep(0.1)  # Simulate I/O
    return {"id": user_id}

# Calling async functions
async def main():
    # Must use 'await' to call async functions
    user = await fetch_user_async(1)
    print(user)

# Run the async code
asyncio.run(main())  # Entry point for async programs
```

### The Three Keywords

```python
# 1. async def - defines a coroutine
async def my_coroutine():
    pass

# 2. await - waits for async operation to complete
async def fetch_data():
    result = await some_async_operation()
    return result

# 3. asyncio.run() - runs the main coroutine
asyncio.run(main())
```

## Running Multiple Operations Concurrently

### asyncio.gather() - Run Multiple Tasks Together

```python
async def fetch_user(user_id: int) -> dict:
    await asyncio.sleep(0.1)
    return {"id": user_id, "name": f"User{user_id}"}

async def fetch_orders(user_id: int) -> list:
    await asyncio.sleep(0.2)
    return [{"id": 1, "item": "Book"}]

async def fetch_profile(user_id: int) -> dict:
    await asyncio.sleep(0.15)
    return {"bio": "Developer"}

# Sequential (slow) - 450ms
async def get_user_data_slow(user_id: int):
    user = await fetch_user(user_id)      # 100ms
    orders = await fetch_orders(user_id)  # 200ms
    profile = await fetch_profile(user_id) # 150ms
    return user, orders, profile

# Concurrent (fast) - 200ms
async def get_user_data_fast(user_id: int):
    # All run at the same time, wait for all to complete
    user, orders, profile = await asyncio.gather(
        fetch_user(user_id),
        fetch_orders(user_id),
        fetch_profile(user_id)
    )
    return user, orders, profile

# With error handling
async def get_user_data_safe(user_id: int):
    results = await asyncio.gather(
        fetch_user(user_id),
        fetch_orders(user_id),
        fetch_profile(user_id),
        return_exceptions=True  # Don't fail all if one fails
    )
    
    # Check for exceptions
    user, orders, profile = results
    if isinstance(orders, Exception):
        orders = []  # Default if failed
    
    return user, orders, profile
```

### asyncio.create_task() - Fire and Continue

```python
async def log_event(event: str):
    await asyncio.sleep(0.5)
    print(f"Logged: {event}")

async def process_request(data: dict):
    # Start logging but don't wait for it
    task = asyncio.create_task(log_event("request_processed"))
    
    # Continue processing immediately
    result = await process_data(data)
    
    # Optionally wait for logging to finish later
    await task
    
    return result

# Multiple background tasks
async def main():
    tasks = [
        asyncio.create_task(process_user(1)),
        asyncio.create_task(process_user(2)),
        asyncio.create_task(process_user(3)),
    ]
    
    # Wait for all tasks
    results = await asyncio.gather(*tasks)
```

### asyncio.wait_for() - Timeout Protection

```python
async def fetch_external_api():
    await asyncio.sleep(10)  # Slow API
    return "data"

async def fetch_with_timeout():
    try:
        # Wait max 5 seconds
        result = await asyncio.wait_for(
            fetch_external_api(),
            timeout=5.0
        )
        return result
    except asyncio.TimeoutError:
        print("API call timed out")
        return None
```

## Async Context Managers

```python
# Async context manager for database connection
class AsyncDatabase:
    async def __aenter__(self):
        self.conn = await create_connection()
        return self.conn
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        await self.conn.close()

# Usage
async def query_database():
    async with AsyncDatabase() as db:
        result = await db.query("SELECT * FROM users")
        return result

# Real-world example with aiohttp
import aiohttp

async def fetch_url(url: str) -> str:
    async with aiohttp.ClientSession() as session:
        async with session.get(url) as response:
            return await response.text()
```

## Async Iteration

```python
# Async generator
async def fetch_users_paginated(page_size: int = 100):
    page = 1
    while True:
        users = await db.query(
            f"SELECT * FROM users LIMIT {page_size} OFFSET {(page-1)*page_size}"
        )
        if not users:
            break
        for user in users:
            yield user
        page += 1

# Usage with async for
async def process_all_users():
    async for user in fetch_users_paginated():
        await process_user(user)

# Async comprehension
async def get_all_user_emails():
    return [
        user["email"] 
        async for user in fetch_users_paginated()
    ]
```

## Common Patterns in Backend Services

### Pattern 1: Parallel API Calls

```python
import aiohttp

async def fetch_json(session: aiohttp.ClientSession, url: str):
    async with session.get(url) as response:
        return await response.json()

async def aggregate_data(user_id: int):
    async with aiohttp.ClientSession() as session:
        # Fetch from multiple services concurrently
        user_data, orders_data, recommendations = await asyncio.gather(
            fetch_json(session, f"https://api.users.com/{user_id}"),
            fetch_json(session, f"https://api.orders.com/user/{user_id}"),
            fetch_json(session, f"https://api.recommend.com/{user_id}")
        )
    
    return {
        "user": user_data,
        "orders": orders_data,
        "recommendations": recommendations
    }
```

### Pattern 2: Database Connection Pooling

```python
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker

# Create async engine with connection pool
engine = create_async_engine(
    "postgresql+asyncpg://user:pass@localhost/db",
    pool_size=20,
    max_overflow=10
)

AsyncSessionLocal = sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False
)

# Dependency for FastAPI
async def get_db():
    async with AsyncSessionLocal() as session:
        yield session

# Usage in endpoint
async def get_user(user_id: int, db: AsyncSession):
    result = await db.execute(
        select(User).where(User.id == user_id)
    )
    return result.scalar_one_or_none()
```

### Pattern 3: Rate Limiting with Semaphore

```python
# Limit concurrent operations
semaphore = asyncio.Semaphore(10)  # Max 10 concurrent

async def fetch_user(user_id: int):
    async with semaphore:  # Acquire semaphore
        await asyncio.sleep(0.1)
        return {"id": user_id}
    # Semaphore auto-released

async def fetch_many_users(user_ids: list[int]):
    tasks = [fetch_user(uid) for uid in user_ids]
    return await asyncio.gather(*tasks)  # Max 10 at a time
```

### Pattern 4: Background Task with Queue

```python
import asyncio
from asyncio import Queue

# Producer-consumer pattern
async def producer(queue: Queue):
    for i in range(100):
        await queue.put(i)
        await asyncio.sleep(0.01)
    await queue.put(None)  # Sentinel to stop

async def consumer(queue: Queue, name: str):
    while True:
        item = await queue.get()
        if item is None:
            queue.put_nowait(None)  # For other consumers
            break
        
        # Process item
        await process_item(item)
        queue.task_done()

async def main():
    queue = Queue(maxsize=10)
    
    # Start producer and multiple consumers
    await asyncio.gather(
        producer(queue),
        consumer(queue, "worker-1"),
        consumer(queue, "worker-2"),
        consumer(queue, "worker-3")
    )
```

## Async vs Threads vs Multiprocessing

```python
# Use asyncio for I/O-bound tasks
async def io_bound():
    async with aiohttp.ClientSession() as session:
        tasks = [
            fetch_url(session, url) 
            for url in urls
        ]
        return await asyncio.gather(*tasks)

# Use threading for blocking I/O (non-async libraries)
import concurrent.futures

def cpu_bound_task(n: int) -> int:
    return sum(i * i for i in range(n))

async def run_blocking_tasks():
    loop = asyncio.get_event_loop()
    with concurrent.futures.ThreadPoolExecutor() as executor:
        tasks = [
            loop.run_in_executor(executor, cpu_bound_task, n)
            for n in [1000000, 2000000, 3000000]
        ]
        return await asyncio.gather(*tasks)

# Use multiprocessing for CPU-intensive tasks
from concurrent.futures import ProcessPoolExecutor

async def run_cpu_intensive():
    loop = asyncio.get_event_loop()
    with ProcessPoolExecutor() as executor:
        result = await loop.run_in_executor(
            executor,
            heavy_computation,
            data
        )
        return result
```

## Common Mistakes to Avoid

```python
# ❌ WRONG: Forgetting await (returns coroutine object, doesn't run)
async def wrong():
    result = fetch_data()  # This doesn't run!
    return result

# ✅ RIGHT
async def correct():
    result = await fetch_data()
    return result

# ❌ WRONG: Using blocking code in async function
import time
async def wrong_sleep():
    time.sleep(5)  # BLOCKS the entire event loop!

# ✅ RIGHT
async def correct_sleep():
    await asyncio.sleep(5)  # Non-blocking

# ❌ WRONG: Mixing sync and async incorrectly
def sync_function():
    return asyncio.run(async_function())  # Don't do this in async context

# ✅ RIGHT: Use run_in_executor for sync code in async
async def correct_way():
    loop = asyncio.get_event_loop()
    result = await loop.run_in_executor(None, sync_blocking_function)

# ❌ WRONG: Not handling exceptions in gather
async def wrong():
    results = await asyncio.gather(task1(), task2())
    # If task1 fails, task2 is cancelled and never returns

# ✅ RIGHT
async def correct():
    results = await asyncio.gather(
        task1(), 
        task2(), 
        return_exceptions=True
    )
    # Check each result for exceptions
```

## Testing Async Code

```python
import pytest

# Mark test as async
@pytest.mark.asyncio
async def test_fetch_user():
    user = await fetch_user(1)
    assert user["id"] == 1

# Test with mocked async function
from unittest.mock import AsyncMock

@pytest.mark.asyncio
async def test_with_mock():
    mock_db = AsyncMock()
    mock_db.query.return_value = {"id": 1, "name": "Alice"}
    
    result = await get_user_data(mock_db, 1)
    assert result["name"] == "Alice"
    mock_db.query.assert_called_once()
```

## Quick Reference

```python
# Run single coroutine
asyncio.run(main())

# Run multiple concurrently
await asyncio.gather(coro1(), coro2(), coro3())

# Create background task
task = asyncio.create_task(background_job())

# Add timeout
await asyncio.wait_for(slow_operation(), timeout=5.0)

# Limit concurrency
semaphore = asyncio.Semaphore(10)
async with semaphore:
    await operation()

# Run blocking code in thread
loop = asyncio.get_event_loop()
await loop.run_in_executor(None, blocking_function)

# Async context manager
async with resource() as r:
    await r.use()

# Async iteration
async for item in async_generator():
    await process(item)
```

## Key Takeaways

1. **Use async for I/O-bound operations** - databases, APIs, file I/O
2. **Always await async functions** - or they won't run
3. **Use asyncio.gather()** for parallel operations
4. **Never use blocking code** in async functions (time.sleep, requests, etc.)
5. **return_exceptions=True** in gather() for resilience
6. **Semaphores** to limit concurrency and avoid overwhelming services

Next: [03-fastapi-modern-backend.md](./03-fastapi-modern-backend.md) - Building production APIs with FastAPI

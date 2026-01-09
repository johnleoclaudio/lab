# Python Essentials for Backend Developers

Time: 1.5 hours | Focus: What you NEED, not everything Python offers

## Type Hints (Use Them Everywhere in 2025)

```python
# Basic types
def process_user(user_id: int, name: str, active: bool = True) -> dict[str, any]:
    return {"id": user_id, "name": name, "active": active}

# Modern syntax (Python 3.10+)
from typing import Optional  # Old way
user_id: int | None = None  # New way, preferred in 2025

# Collections
def get_users(ids: list[int]) -> dict[int, str]:
    return {1: "Alice", 2: "Bob"}

# Type aliases for complex types
UserId = int
UserData = dict[str, str | int]

def create_user(user_id: UserId) -> UserData:
    return {"id": user_id, "name": "test"}
```

**Why it matters**: FastAPI uses type hints for validation, Pydantic uses them for models, IDEs love them.

## Data Structures You'll Actually Use

### Dictionaries (Your Best Friend)

```python
# Creation
user = {"id": 1, "name": "Alice", "email": "alice@example.com"}

# Get with default (avoid KeyError)
role = user.get("role", "guest")  # Returns "guest" if key doesn't exist

# Merge dictionaries (Python 3.9+)
defaults = {"role": "user", "active": True}
user = defaults | {"id": 1, "name": "Alice"}  # Prefer this over update()

# Dictionary comprehension
squared = {x: x**2 for x in range(5)}  # {0: 0, 1: 1, 2: 4, 3: 9, 4: 16}

# Unpacking
def create_user(name: str, email: str, role: str = "user"):
    pass

user_data = {"name": "Alice", "email": "alice@example.com"}
create_user(**user_data)  # Unpacks dict as keyword arguments
```

### Lists and List Comprehensions

```python
# List comprehension (prefer over map/filter)
ids = [1, 2, 3, 4, 5]
even_ids = [x for x in ids if x % 2 == 0]  # [2, 4]

# With transformation
squared_evens = [x**2 for x in ids if x % 2 == 0]  # [4, 16]

# Multiple operations
users = [{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}]
names = [u["name"].upper() for u in users]  # ["ALICE", "BOB"]

# Unpacking
first, *rest = [1, 2, 3, 4]  # first=1, rest=[2, 3, 4]
first, *middle, last = [1, 2, 3, 4]  # first=1, middle=[2, 3], last=4
```

### Sets (For Deduplication and Fast Lookups)

```python
# Remove duplicates
ids = [1, 2, 2, 3, 3, 4]
unique_ids = list(set(ids))  # [1, 2, 3, 4]

# Fast membership testing (O(1) vs O(n) for lists)
allowed_roles = {"admin", "moderator", "user"}
if user_role in allowed_roles:  # Much faster than list
    pass

# Set operations
active_users = {1, 2, 3, 4}
premium_users = {3, 4, 5, 6}
premium_active = active_users & premium_users  # Intersection: {3, 4}
all_special = active_users | premium_users  # Union: {1, 2, 3, 4, 5, 6}
```

## Functions and Common Patterns

### Default Arguments (Be Careful!)

```python
# WRONG - mutable default argument bug
def add_user(users=[]):  # DON'T DO THIS
    users.append("new")
    return users

# RIGHT
def add_user(users: list[str] | None = None) -> list[str]:
    if users is None:
        users = []
    users.append("new")
    return users

# Even better with walrus operator (Python 3.8+)
def add_user(users: list[str] | None = None) -> list[str]:
    return (users if users is not None else []) + ["new"]
```

### Args and Kwargs (For Flexibility)

```python
# *args for variable positional arguments
def sum_all(*args: int) -> int:
    return sum(args)

sum_all(1, 2, 3, 4)  # 10

# **kwargs for variable keyword arguments
def create_user(name: str, **kwargs) -> dict:
    user = {"name": name}
    user.update(kwargs)
    return user

create_user("Alice", role="admin", active=True)

# Combine everything (order matters: positional, *args, keyword, **kwargs)
def complex_function(required: str, *args, default: int = 0, **kwargs):
    pass
```

### Lambda Functions (Keep Them Simple)

```python
# Good use: simple sorting
users = [{"name": "Charlie", "age": 30}, {"name": "Alice", "age": 25}]
sorted_users = sorted(users, key=lambda u: u["age"])

# Good use: simple filtering
numbers = [1, 2, 3, 4, 5]
evens = list(filter(lambda x: x % 2 == 0, numbers))

# Bad use: complex logic (use a proper function instead)
# Don't do this: lambda x: x if x > 0 else -x if x < -10 else 0
```

## Error Handling (Do It Right)

```python
# Specific exceptions
def get_user(user_id: int) -> dict:
    try:
        user = database.get(user_id)
        return user
    except KeyError:
        raise ValueError(f"User {user_id} not found")
    except ConnectionError as e:
        # Log and re-raise or handle
        logger.error(f"Database connection failed: {e}")
        raise

# Context manager for cleanup (auto-closes resources)
with open("file.txt") as f:
    content = f.read()  # File auto-closes even if error occurs

# Custom context manager (useful for database sessions)
from contextlib import contextmanager

@contextmanager
def database_session():
    session = create_session()
    try:
        yield session
        session.commit()
    except Exception:
        session.rollback()
        raise
    finally:
        session.close()

# Usage
with database_session() as session:
    session.add(user)
```

## String Operations

```python
# F-strings (always use these, not % or .format())
name = "Alice"
age = 30
message = f"User {name} is {age} years old"

# With expressions
message = f"Next year: {age + 1}"

# Format numbers
price = 19.99
formatted = f"Price: ${price:.2f}"  # "Price: $19.99"

# Multi-line strings
query = """
    SELECT id, name, email
    FROM users
    WHERE active = true
    AND role = 'admin'
"""

# String methods you'll use
email = "  ALICE@EXAMPLE.COM  "
clean = email.strip().lower()  # "alice@example.com"

path = "/api/v1/users"
parts = path.split("/")  # ["", "api", "v1", "users"]
```

## Comprehensions (Write Less, Do More)

```python
# List comprehension
squares = [x**2 for x in range(10)]

# Dict comprehension
user_map = {u["id"]: u["name"] for u in users}

# Set comprehension
unique_domains = {email.split("@")[1] for email in email_list}

# Nested comprehension (flatten list)
matrix = [[1, 2], [3, 4], [5, 6]]
flat = [num for row in matrix for num in row]  # [1, 2, 3, 4, 5, 6]

# With conditions
adults = [u for u in users if u["age"] >= 18]

# Complex filtering
active_admins = [
    u["email"]
    for u in users
    if u["active"] and u["role"] == "admin"
]
```

## Working with None (Avoid Null Pointer Errors)

```python
# Check for None explicitly
user_id: int | None = get_user_id()

if user_id is not None:  # Use 'is', not '=='
    process(user_id)

# Default values
name = user_name or "Guest"  # Watch out: 0, "", [] also trigger default
name = user_name if user_name is not None else "Guest"  # Better

# Walrus operator for cleaner code (Python 3.8+)
if (user := get_user()) is not None:
    print(user["name"])  # user is available in this scope

# Optional chaining equivalent (use get() and defaults)
email = user.get("contact", {}).get("email", "no-email@example.com")
```

## Imports and Modules

```python
# Absolute imports (preferred)
from myapp.services.user import UserService
from myapp.models import User, Role

# Relative imports (only within packages)
from .models import User  # Same directory
from ..services import UserService  # Parent directory

# Import aliases
import pandas as pd
from datetime import datetime as dt

# Avoid star imports
from module import *  # DON'T - hard to track what's imported

# Conditional imports (for optional dependencies)
try:
    import redis
except ImportError:
    redis = None
```

## Useful Built-in Functions

```python
# enumerate (get index + value)
for i, user in enumerate(users):
    print(f"{i}: {user['name']}")

# zip (combine iterables)
names = ["Alice", "Bob"]
ages = [25, 30]
users = [{"name": n, "age": a} for n, a in zip(names, ages)]

# any/all (check conditions)
has_admin = any(u["role"] == "admin" for u in users)
all_active = all(u["active"] for u in users)

# sorted (don't modify original)
sorted_users = sorted(users, key=lambda u: u["age"], reverse=True)

# min/max with key
youngest = min(users, key=lambda u: u["age"])
oldest = max(users, key=lambda u: u["age"])
```

## Practice Exercise

Create a function that processes a list of user dictionaries:

```python
def process_users(users: list[dict]) -> dict:
    """
    Given users with keys: id, name, email, age, role, active
    Return:
    {
        "active_count": int,
        "adult_admins": list[str] (emails),
        "average_age": float,
        "role_distribution": dict[str, int]
    }
    """
    pass

# Example usage
users = [
    {"id": 1, "name": "Alice", "email": "alice@ex.com", "age": 25, "role": "admin", "active": True},
    {"id": 2, "name": "Bob", "email": "bob@ex.com", "age": 17, "role": "user", "active": True},
    {"id": 3, "name": "Charlie", "email": "charlie@ex.com", "age": 30, "role": "admin", "active": False},
]

result = process_users(users)
```

<details>
<summary>Solution</summary>

```python
def process_users(users: list[dict]) -> dict:
    return {
        "active_count": sum(1 for u in users if u["active"]),
        "adult_admins": [
            u["email"] 
            for u in users 
            if u["age"] >= 18 and u["role"] == "admin"
        ],
        "average_age": sum(u["age"] for u in users) / len(users) if users else 0,
        "role_distribution": {
            role: sum(1 for u in users if u["role"] == role)
            for role in {u["role"] for u in users}
        }
    }
```

</details>

## Key Takeaways

1. **Type hints everywhere** - not optional in modern Python backend
2. **Comprehensions over loops** - more Pythonic, often faster
3. **Use `is None` not `== None`** - identity vs equality
4. **Context managers** for resource management
5. **F-strings** for all string formatting
6. **Dict operations** are your daily bread

Next: [02-async-concurrency.md](./02-async-concurrency.md) - Understanding async/await for backend services
